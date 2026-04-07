package i2ptui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-i2p/i2ptui/rpc"
)

type controlAction struct {
	name    string
	desc    string
	execute func() (string, error)
}

type controlActionMsg struct {
	result string
}

type controlModel struct {
	actions    []controlAction
	cursor     int
	confirming bool
}

func newControlModel() controlModel {
	return controlModel{
		actions: []controlAction{
			{
				name: "Restart (graceful)",
				desc: "The router will restart in ~11 minutes",
				execute: func() (string, error) {
					return rpc.RestartGraceful()
				},
			},
			{
				name: "Restart (immediate)",
				desc: "The router will restart immediately",
				execute: func() (string, error) {
					return rpc.Restart()
				},
			},
			{
				name: "Shutdown (graceful)",
				desc: "The router will shut down in ~11 minutes",
				execute: func() (string, error) {
					return rpc.ShutdownGraceful()
				},
			},
			{
				name: "Shutdown (immediate)",
				desc: "The router will shut down immediately",
				execute: func() (string, error) {
					return rpc.Shutdown()
				},
			},
			{
				name: "Check for Updates",
				desc: "Check if router updates are available",
				execute: func() (string, error) {
					return rpc.FindUpdates()
				},
			},
		},
	}
}

// Update handles key events for the control tab.
func (m controlModel) Update(msg tea.Msg) (controlModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.confirming {
			switch msg.String() {
			case "y", "enter":
				m.confirming = false
				action := m.actions[m.cursor]
				return m, m.executeAction(action)
			case "n", "esc":
				m.confirming = false
				return m, nil
			}
			return m, nil
		}
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.actions)-1 {
				m.cursor++
			}
		case "enter":
			m.confirming = true
		}
	}
	return m, nil
}

func (m controlModel) executeAction(a controlAction) tea.Cmd {
	return func() tea.Msg {
		result, err := a.execute()
		if err != nil {
			return controlActionMsg{
				result: fmt.Sprintf("Error: %s: %v", a.name, err),
			}
		}
		return controlActionMsg{
			result: fmt.Sprintf("%s: %s", a.name, result),
		}
	}
}

// View renders the control tab.
func (m controlModel) View(width int) string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render("  Router Control"))
	b.WriteString("\n\n")

	for i, a := range m.actions {
		cursor := "  "
		style := lipgloss.NewStyle()
		if i == m.cursor {
			cursor = "> "
			style = selectedStyle
		}
		b.WriteString(style.Render(fmt.Sprintf("  %s%s", cursor, a.name)))
		b.WriteString("\n")
	}

	b.WriteString("\n  Press Enter to select, q to quit\n")

	if m.confirming {
		a := m.actions[m.cursor]
		box := fmt.Sprintf(
			"%s\n%s\n\n  [y] Yes    [n] No",
			selectedStyle.Render(a.name+"?"),
			a.desc,
		)
		b.WriteString("\n")
		b.WriteString(confirmBoxStyle.Render(box))
		b.WriteString("\n")
	}

	return b.String()
}
