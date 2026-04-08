package i2ptui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestApplyDarkTheme(t *testing.T) {
	applyTheme(DarkTheme)
	// Verify it doesn't panic — styles are global, so just ensure no error.
}

func TestApplyLightTheme(t *testing.T) {
	applyTheme(LightTheme)
	applyTheme(DarkTheme) // Reset to default.
}

func TestWithThemeOption(t *testing.T) {
	_ = New(WithTheme(LightTheme))
	applyTheme(DarkTheme) // Reset.
}

func TestTabFromClick(t *testing.T) {
	m := New()
	m.width = 80

	// Click at x=0 should be first tab (Dashboard).
	tab := m.tabFromClick(0)
	if tab != tabDashboard {
		t.Errorf("expected tabDashboard at x=0, got %d", tab)
	}

	// Far right should return current tab (no tab there).
	tab = m.tabFromClick(200)
	if tab != m.activeTab {
		t.Errorf("expected current tab for out-of-range click, got %d", tab)
	}
}

func TestMouseTabSwitch(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	// Simulate mouse click on Y=0 (tab bar).
	mouseMsg := tea.MouseMsg{
		X:      0,
		Y:      0,
		Action: tea.MouseActionRelease,
		Button: tea.MouseButtonLeft,
	}
	updated, _ := m.Update(mouseMsg)
	m = updated.(Model)
	// Should be on Dashboard (x=0).
	if m.activeTab != tabDashboard {
		t.Errorf("expected tabDashboard after click at 0,0")
	}
}
