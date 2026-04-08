package i2ptui

import (
	"strings"
	"testing"
	"time"

	"github.com/go-i2p/i2ptui/rpc"
)

func TestNotifyStatusChange(t *testing.T) {
	m := newNotifyModel()

	// First snapshot sets baseline — no notification.
	snap1 := rpc.RouterSnapshot{
		Status:    "OK",
		NetStatus: "OK",
		FetchedAt: time.Now(),
	}
	m.CheckChanges(rpc.RouterSnapshot{}, snap1)
	if m.HasNotifications() {
		t.Error("should not notify on first snapshot")
	}

	// Status changes → notification.
	snap2 := rpc.RouterSnapshot{
		Status:    "Testing",
		NetStatus: "OK",
		FetchedAt: time.Now(),
	}
	m.CheckChanges(snap1, snap2)
	if !m.HasNotifications() {
		t.Error("expected notification on status change")
	}
	if !strings.Contains(m.View(), "OK → Testing") {
		t.Error("notification should mention status change")
	}
}

func TestNotifyNetStatusChange(t *testing.T) {
	m := newNotifyModel()

	snap1 := rpc.RouterSnapshot{Status: "OK", NetStatus: "OK"}
	m.CheckChanges(rpc.RouterSnapshot{}, snap1)

	snap2 := rpc.RouterSnapshot{Status: "OK", NetStatus: "Firewalled"}
	m.CheckChanges(snap1, snap2)
	if !m.HasNotifications() {
		t.Error("expected notification on net status change")
	}
}

func TestNotifyReseed(t *testing.T) {
	m := newNotifyModel()

	snap1 := rpc.RouterSnapshot{Status: "OK", Reseeding: false}
	m.CheckChanges(rpc.RouterSnapshot{}, snap1)

	snap2 := rpc.RouterSnapshot{Status: "OK", Reseeding: true}
	m.CheckChanges(snap1, snap2)
	if !m.HasNotifications() {
		t.Error("expected notification on reseed trigger")
	}
}

func TestNotifyDismiss(t *testing.T) {
	m := newNotifyModel()
	m.add("test notification")
	if !m.HasNotifications() {
		t.Error("expected notification")
	}
	m.Dismiss()
	if m.HasNotifications() {
		t.Error("expected no notifications after dismiss")
	}
}

func TestNotifyViewEmpty(t *testing.T) {
	m := newNotifyModel()
	if v := m.View(); v != "" {
		t.Errorf("expected empty view, got %q", v)
	}
}

func TestNotifyDismissEmpty(t *testing.T) {
	m := newNotifyModel()
	m.Dismiss() // should not panic
	if m.HasNotifications() {
		t.Error("expected no notifications")
	}
}
