package i2ptui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-i2p/i2ptui/rpc"
)

func TestGraphToggle(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 40
	if m.showGraphs {
		t.Error("graphs should be hidden by default")
	}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = updated.(Model)
	if !m.showGraphs {
		t.Error("graphs should be visible after pressing g")
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = updated.(Model)
	if m.showGraphs {
		t.Error("graphs should be hidden after pressing g again")
	}
}

func TestHistoryRecording(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 40
	snap := rpc.RouterSnapshot{
		IncomingBW:           51200,
		OutgoingBW:           12800,
		ParticipatingTunnels: 350,
		ExplBuildSuccessPct:  72,
		FetchedAt:            time.Now(),
	}
	updated, _ := m.Update(snap)
	m = updated.(Model)

	if m.inBWHistory.Len() != 1 {
		t.Errorf("expected 1 inBW sample, got %d", m.inBWHistory.Len())
	}
	if m.outBWHistory.Len() != 1 {
		t.Errorf("expected 1 outBW sample, got %d", m.outBWHistory.Len())
	}
	if m.tunnelHistory.Len() != 1 {
		t.Errorf("expected 1 tunnel sample, got %d", m.tunnelHistory.Len())
	}
	if m.buildSuccHistory.Len() != 1 {
		t.Errorf("expected 1 build sample, got %d", m.buildSuccHistory.Len())
	}
}

func TestGraphViewDashboard(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 40
	m.showGraphs = true
	m.loading = false
	snap := rpc.RouterSnapshot{
		Status:               "OK",
		IncomingBW:           51200,
		OutgoingBW:           12800,
		ParticipatingTunnels: 350,
		FetchedAt:            time.Now(),
	}
	m.snapshot = snap
	m.overview = m.overview.SetSnapshot(snap)
	m.inBWHistory.Add(time.Now(), 51200)
	m.outBWHistory.Add(time.Now(), 12800)
	m.tunnelHistory.Add(time.Now(), 350)

	v := m.View()
	if !strings.Contains(v, "Graphs") {
		t.Error("expected Graphs section when showGraphs is true")
	}
	if !strings.Contains(v, "Incoming BW") {
		t.Error("expected Incoming BW sparkline label")
	}
}

func TestGraphViewStats(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 40
	m.activeTab = tabStats
	m.showGraphs = true
	m.loading = false
	snap := rpc.RouterSnapshot{
		ExplBuildSuccessPct: 72,
		ExplBuildRejectPct:  18,
		ExplBuildExpirePct:  10,
		FetchedAt:           time.Now(),
	}
	m.snapshot = snap
	m.stats = m.stats.SetSnapshot(snap)

	v := m.View()
	if !strings.Contains(v, "Build Success Rate") {
		t.Error("expected Build Success Rate section in stats with graphs on")
	}
}

func TestNoHistoryOnError(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 40
	snap := rpc.RouterSnapshot{
		Err:       errTest,
		FetchedAt: time.Now(),
	}
	updated, _ := m.Update(snap)
	m = updated.(Model)
	if m.inBWHistory.Len() != 0 {
		t.Error("should not record history on error snapshot")
	}
}
