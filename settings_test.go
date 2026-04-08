package i2ptui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-i2p/i2ptui/rpc"
)

func TestSettingsModelInit(t *testing.T) {
	m := newSettingsModel()
	if m.loaded {
		t.Error("settings should not be loaded initially")
	}
	if len(m.fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(m.fields))
	}
}

func TestSettingsApply(t *testing.T) {
	m := newSettingsModel()
	s := rpc.RouterSettings{
		BWIn:    "512",
		BWOut:   "128",
		BWShare: "80",
		Upnp:    "true",
	}
	msg := settingsMsg{settings: s}
	m, _ = m.Update(msg)
	if !m.loaded {
		t.Error("expected loaded=true after settings message")
	}
	if m.fields[0].input.Value() != "512" {
		t.Errorf("expected BWIn=512, got %s", m.fields[0].input.Value())
	}
	if m.fields[1].input.Value() != "128" {
		t.Errorf("expected BWOut=128, got %s", m.fields[1].input.Value())
	}
	if m.fields[2].input.Value() != "80" {
		t.Errorf("expected BWShare=80, got %s", m.fields[2].input.Value())
	}
}

func TestSettingsNavigation(t *testing.T) {
	m := newSettingsModel()
	m.loaded = true

	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1 after down, got %d", m.cursor)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0 after up, got %d", m.cursor)
	}
}

func TestSettingsConfirmation(t *testing.T) {
	m := newSettingsModel()
	m.loaded = true

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !m.confirming {
		t.Error("expected confirming=true after Enter")
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m.confirming {
		t.Error("expected confirming=false after Esc")
	}
}

func TestSettingsView(t *testing.T) {
	m := newSettingsModel()
	m.loaded = true
	v := m.View(80)

	checks := []string{"Router Settings", "Incoming BW", "Outgoing BW", "Share %"}
	for _, c := range checks {
		if !strings.Contains(v, c) {
			t.Errorf("settings view missing %q", c)
		}
	}
}

func TestSettingsViewLoading(t *testing.T) {
	m := newSettingsModel()
	v := m.View(80)
	if !strings.Contains(v, "Loading") {
		t.Error("expected Loading text when settings not loaded")
	}
}

func TestSettingsRestartIndicator(t *testing.T) {
	m := newSettingsModel()
	m.loaded = true
	m.needRestart = true
	v := m.View(80)
	if !strings.Contains(v, "Restart required") {
		t.Error("expected restart indicator when needRestart=true")
	}
}

func TestSettingsConfirmView(t *testing.T) {
	m := newSettingsModel()
	m.loaded = true
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	v := m.View(80)
	if !strings.Contains(v, "Apply settings?") {
		t.Error("expected confirm dialog")
	}
}

func TestSettingsSaved(t *testing.T) {
	m := newSettingsModel()
	m.loaded = true
	m.saving = true

	m, _ = m.Update(settingsSavedMsg{})
	if m.saving {
		t.Error("expected saving=false after save complete")
	}
	if m.lastOK != "Settings saved" {
		t.Errorf("expected success message, got %q", m.lastOK)
	}
	if !m.needRestart {
		t.Error("expected needRestart=true after save")
	}
}
