package rpc

import (
	"testing"
	"time"
)

func TestRouterSnapshotUptimeDuration(t *testing.T) {
	s := RouterSnapshot{Uptime: 3600000}
	d := s.UptimeDuration()
	if d != time.Hour {
		t.Errorf("expected 1h, got %v", d)
	}
}

func TestRouterSnapshotZeroUptime(t *testing.T) {
	s := RouterSnapshot{}
	d := s.UptimeDuration()
	if d != 0 {
		t.Errorf("expected 0, got %v", d)
	}
}
