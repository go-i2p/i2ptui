package i2ptui

import (
	"testing"
	"time"
)

func TestHistoryBufferAdd(t *testing.T) {
	h := newHistoryBuffer(5 * time.Minute)
	now := time.Now()
	h.Add(now, 100)
	h.Add(now.Add(time.Second), 200)

	vals := h.Values()
	if len(vals) != 2 {
		t.Fatalf("expected 2 values, got %d", len(vals))
	}
	if vals[0] != 100 || vals[1] != 200 {
		t.Errorf("unexpected values: %v", vals)
	}
}

func TestHistoryBufferExpiry(t *testing.T) {
	h := newHistoryBuffer(2 * time.Second)
	base := time.Now().Add(-5 * time.Second)

	h.Add(base, 10)                            // expired
	h.Add(base.Add(1*time.Second), 20)         // expired
	h.Add(base.Add(4*time.Second), 30)         // recent
	h.Add(base.Add(4500*time.Millisecond), 40) // recent

	vals := h.Values()
	if len(vals) != 2 {
		t.Fatalf("expected 2 values after expiry, got %d", len(vals))
	}
	if vals[0] != 30 || vals[1] != 40 {
		t.Errorf("unexpected values after expiry: %v", vals)
	}
}

func TestHistoryBufferEmpty(t *testing.T) {
	h := newHistoryBuffer(time.Minute)
	vals := h.Values()
	if len(vals) != 0 {
		t.Errorf("expected empty, got %v", vals)
	}
}

func TestHistoryBufferLen(t *testing.T) {
	h := newHistoryBuffer(time.Minute)
	if h.Len() != 0 {
		t.Errorf("expected 0, got %d", h.Len())
	}
	h.Add(time.Now(), 42)
	if h.Len() != 1 {
		t.Errorf("expected 1, got %d", h.Len())
	}
}
