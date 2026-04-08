package i2ptui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-i2p/i2ptui/rpc"
)

// settingsMsg carries fetched settings from the RPC layer.
type settingsMsg struct {
	settings rpc.RouterSettings
}

// settingsSavedMsg signals settings were saved.
type settingsSavedMsg struct {
	err error
}

type settingsField struct {
	label string
	key   string
	input textinput.Model
}

type settingsModel struct {
	fields      []settingsField
	cursor      int
	loaded      bool
	confirming  bool
	saving      bool
	lastErr     string
	lastOK      string
	needRestart bool
}

// newSettingsModel returns a settingsModel with default fields.
func newSettingsModel() settingsModel {
	fields := []settingsField{
		newField("Incoming BW (KB/s)", "i2p.router.net.bw.in"),
		newField("Outgoing BW (KB/s)", "i2p.router.net.bw.out"),
		newField("Share %", "i2p.router.net.bw.share"),
	}
	fields[0].input.Focus()
	return settingsModel{fields: fields}
}

// newField creates a settingsField with the given label and RPC key.
func newField(label, key string) settingsField {
	ti := textinput.New()
	ti.Placeholder = "..."
	ti.CharLimit = 10
	ti.Width = 12
	return settingsField{label: label, key: key, input: ti}
}

// fetchSettingsCmd is a tea.Cmd that reads router settings via RPC.
func fetchSettingsCmd() tea.Msg {
	return settingsMsg{settings: rpc.ReadSettings()}
}

// saveCmd returns a tea.Cmd that writes the current field values to the router.
func (m settingsModel) saveCmd() tea.Cmd {
	return func() tea.Msg {
		for _, f := range m.fields {
			val := strings.TrimSpace(f.input.Value())
			if val == "" {
				continue
			}
			if err := rpc.WriteSetting(f.key, val); err != nil {
				return settingsSavedMsg{err: err}
			}
		}
		return settingsSavedMsg{}
	}
}

// Update handles key events for the settings tab.
func (m settingsModel) Update(msg tea.Msg) (settingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case settingsMsg:
		m.loaded = true
		m.applySettings(msg.settings)
		return m, nil

	case settingsSavedMsg:
		m.saving = false
		if msg.err != nil {
			m.lastErr = msg.err.Error()
			m.lastOK = ""
		} else {
			m.lastOK = "Settings saved"
			m.lastErr = ""
			m.needRestart = true
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m.updateInputs(msg)
}

// handleKey dispatches key events for navigation, editing, and saving.
func (m settingsModel) handleKey(msg tea.KeyMsg) (settingsModel, tea.Cmd) {
	if m.confirming {
		return m.handleConfirm(msg)
	}
	switch msg.String() {
	case "tab", "down", "j":
		m.fields[m.cursor].input.Blur()
		m.cursor = (m.cursor + 1) % len(m.fields)
		m.fields[m.cursor].input.Focus()
	case "shift+tab", "up", "k":
		m.fields[m.cursor].input.Blur()
		m.cursor = (m.cursor - 1 + len(m.fields)) % len(m.fields)
		m.fields[m.cursor].input.Focus()
	case "enter":
		m.confirming = true
	}
	return m, nil
}

// handleConfirm handles key events within the save confirmation dialog.
func (m settingsModel) handleConfirm(msg tea.KeyMsg) (settingsModel, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		m.confirming = false
		m.saving = true
		return m, m.saveCmd()
	case "n", "esc":
		m.confirming = false
	}
	return m, nil
}

// updateInputs forwards messages to all text input fields.
func (m settingsModel) updateInputs(msg tea.Msg) (settingsModel, tea.Cmd) {
	var cmds []tea.Cmd
	for i := range m.fields {
		var cmd tea.Cmd
		m.fields[i].input, cmd = m.fields[i].input.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// applySettings populates the text fields from fetched RouterSettings.
func (m *settingsModel) applySettings(s rpc.RouterSettings) {
	values := []string{s.BWIn, s.BWOut, s.BWShare}
	for i, v := range values {
		if i < len(m.fields) && v != "N/A" {
			m.fields[i].input.SetValue(v)
		}
	}
}

// View renders the settings tab.
func (m settingsModel) View(width int) string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render("  Router Settings"))
	b.WriteString("\n\n")

	if !m.loaded {
		b.WriteString("  Loading settings...\n")
		return b.String()
	}

	for i, f := range m.fields {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		b.WriteString(fmt.Sprintf("  %s%s  %s\n",
			cursor,
			labelStyle.Render(f.label),
			f.input.View(),
		))
	}

	b.WriteString("\n  Press Enter to save, Tab to switch fields\n")

	if m.needRestart {
		b.WriteString(restartStyle.Render("  ⚠ Restart required for changes to take effect"))
		b.WriteString("\n")
	}
	b.WriteString(m.renderStatus())
	b.WriteString(m.renderConfirm())

	return b.String()
}

// renderStatus returns the save-in-progress, error, or success message.
func (m settingsModel) renderStatus() string {
	if m.saving {
		return "  Saving...\n"
	}
	if m.lastErr != "" {
		return errorStyle.Render(fmt.Sprintf("  Error: %s", m.lastErr)) + "\n"
	}
	if m.lastOK != "" {
		return statusStyle.Render(fmt.Sprintf("  %s", m.lastOK)) + "\n"
	}
	return ""
}

// renderConfirm returns the confirmation dialog, or an empty string if hidden.
func (m settingsModel) renderConfirm() string {
	if !m.confirming {
		return ""
	}
	box := fmt.Sprintf(
		"%s\n%s\n\n  [y] Yes    [n] No",
		selectedStyle.Render("Apply settings?"),
		"This may require a router restart",
	)
	return "\n" + confirmBoxStyle.Render(box) + "\n"
}

var restartStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("214")).
	Bold(true)
