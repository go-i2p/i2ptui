package i2ptui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-i2p/i2ptui/rpc"
)

func TestNewDefaults(t *testing.T) {
	m := New()
	if m.host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", m.host)
	}
	if m.port != "7657" {
		t.Errorf("expected port 7657, got %s", m.port)
	}
	if m.path != "jsonrpc" {
		t.Errorf("expected path jsonrpc, got %s", m.path)
	}
	if m.password != "itoopie" {
		t.Errorf("expected password itoopie, got %s", m.password)
	}
	if m.interval != 5*time.Second {
		t.Errorf("expected interval 5s, got %v", m.interval)
	}
}

func TestNewWithOptions(t *testing.T) {
	m := New(
		WithHost("10.0.0.1"),
		WithPort("7658"),
		WithPath("rpc"),
		WithPassword("secret"),
		WithCert("/tmp/cert.pem"),
		WithInterval(10*time.Second),
	)
	if m.host != "10.0.0.1" {
		t.Errorf("expected host 10.0.0.1, got %s", m.host)
	}
	if m.port != "7658" {
		t.Errorf("expected port 7658, got %s", m.port)
	}
	if m.path != "rpc" {
		t.Errorf("expected path rpc, got %s", m.path)
	}
	if m.password != "secret" {
		t.Errorf("expected password secret, got %s", m.password)
	}
	if m.cert != "/tmp/cert.pem" {
		t.Errorf("expected cert /tmp/cert.pem, got %s", m.cert)
	}
	if m.interval != 10*time.Second {
		t.Errorf("expected interval 10s, got %v", m.interval)
	}
}

func TestTabSwitching(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24
	if m.activeTab != tabDashboard {
		t.Errorf("expected initial tab dashboard, got %d", m.activeTab)
	}

	// Tab forward
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	m = updated.(Model)
	if m.activeTab != tabStats {
		t.Errorf("expected tab stats after '2', got %d", m.activeTab)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	m = updated.(Model)
	if m.activeTab != tabPeers {
		t.Errorf("expected tab peers after '3', got %d", m.activeTab)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	m = updated.(Model)
	if m.activeTab != tabControl {
		t.Errorf("expected tab control after '4', got %d", m.activeTab)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)
	if m.activeTab != tabDashboard {
		t.Errorf("expected tab dashboard after '1', got %d", m.activeTab)
	}

	// Tab key cycles
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.activeTab != tabStats {
		t.Errorf("expected tab stats after Tab, got %d", m.activeTab)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = updated.(Model)
	if m.activeTab != tabDashboard {
		t.Errorf("expected tab dashboard after Shift+Tab, got %d", m.activeTab)
	}
}

func TestSnapshotUpdate(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24
	m.loading = true
	snap := rpc.RouterSnapshot{
		Status:    "OK",
		NetStatus: "OK",
		Version:   "0.9.62",
		Uptime:    3600000,
		FetchedAt: time.Now(),
	}
	updated, _ := m.Update(snap)
	m = updated.(Model)
	if m.loading {
		t.Error("expected loading to be false after snapshot")
	}
	if m.snapshot.Status != "OK" {
		t.Errorf("expected status OK, got %s", m.snapshot.Status)
	}
	if m.err != nil {
		t.Errorf("expected no error, got %v", m.err)
	}
}

var errTest = &testError{}

type testError struct{}

func (e *testError) Error() string { return "test error" }

func TestSnapshotWithError(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24
	snap := rpc.RouterSnapshot{
		Err:       errTest,
		FetchedAt: time.Now(),
	}
	updated, _ := m.Update(snap)
	m = updated.(Model)
	if m.err == nil {
		t.Error("expected error to be set")
	}
}

func TestViewInitializing(t *testing.T) {
	m := New()
	// width=0 means no WindowSizeMsg yet
	v := m.View()
	if !strings.Contains(v, "Initializing") {
		t.Error("expected Initializing message when width is 0")
	}
}

func TestViewWithSnapshot(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24
	m.loading = false
	m.snapshot = rpc.RouterSnapshot{
		Status:    "OK",
		NetStatus: "Testing",
		Version:   "0.9.62",
		FetchedAt: time.Now(),
	}
	m.overview = m.overview.SetSnapshot(m.snapshot)
	v := m.View()
	if !strings.Contains(v, "OK") {
		t.Error("expected view to contain router status")
	}
	if !strings.Contains(v, "Dashboard") {
		t.Error("expected view to contain Dashboard tab")
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := New()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)
	if m.width != 120 || m.height != 40 {
		t.Errorf("expected 120x40, got %dx%d", m.width, m.height)
	}
}
