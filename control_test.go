package i2ptui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestControlNavigation(t *testing.T) {
	m := newControlModel()

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

	// Can't go above 0
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Errorf("expected cursor clamped at 0, got %d", m.cursor)
	}
}

func TestControlConfirmation(t *testing.T) {
	m := newControlModel()

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !m.confirming {
		t.Error("expected confirming to be true after Enter")
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m.confirming {
		t.Error("expected confirming to be false after Esc")
	}
}

func TestControlView(t *testing.T) {
	m := newControlModel()
	v := m.View(80)

	if !strings.Contains(v, "Restart (graceful)") {
		t.Error("control view missing Restart (graceful)")
	}
	if !strings.Contains(v, "Shutdown (immediate)") {
		t.Error("control view missing Shutdown (immediate)")
	}
	if !strings.Contains(v, "Check for Updates") {
		t.Error("control view missing Check for Updates")
	}
}

func TestControlConfirmView(t *testing.T) {
	m := newControlModel()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	v := m.View(80)
	if !strings.Contains(v, "Yes") {
		t.Error("confirm dialog missing Yes")
	}
	if !strings.Contains(v, "No") {
		t.Error("confirm dialog missing No")
	}
}
