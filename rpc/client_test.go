package rpc

import (
	"errors"
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

func TestFetchStringSuccess(t *testing.T) {
	fn := func() (string, error) { return "OK", nil }
	if v := fetchString(fn); v != "OK" {
		t.Errorf("expected OK, got %s", v)
	}
}

func TestFetchStringError(t *testing.T) {
	fn := func() (string, error) { return "", errors.New("fail") }
	if v := fetchString(fn); v != "N/A" {
		t.Errorf("expected N/A, got %s", v)
	}
}

func TestFetchIntSuccess(t *testing.T) {
	fn := func() (int, error) { return 42, nil }
	if v := fetchInt(fn); v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestFetchIntError(t *testing.T) {
	fn := func() (int, error) { return 0, errors.New("fail") }
	if v := fetchInt(fn); v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}

func TestFetchInt64Success(t *testing.T) {
	fn := func() (int64, error) { return 100, nil }
	if v := fetchInt64(fn); v != 100 {
		t.Errorf("expected 100, got %d", v)
	}
}

func TestFetchInt64Error(t *testing.T) {
	fn := func() (int64, error) { return 0, errors.New("fail") }
	if v := fetchInt64(fn); v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}

func TestFetchFloat64Success(t *testing.T) {
	fn := func() (float64, error) { return 3.14, nil }
	if v := fetchFloat64(fn); v != 3.14 {
		t.Errorf("expected 3.14, got %f", v)
	}
}

func TestFetchFloat64Error(t *testing.T) {
	fn := func() (float64, error) { return 0, errors.New("fail") }
	if v := fetchFloat64(fn); v != 0 {
		t.Errorf("expected 0, got %f", v)
	}
}

func TestFetchBoolSuccess(t *testing.T) {
	fn := func() (bool, error) { return true, nil }
	if v := fetchBool(fn); !v {
		t.Error("expected true")
	}
}

func TestFetchBoolError(t *testing.T) {
	fn := func() (bool, error) { return false, errors.New("fail") }
	if v := fetchBool(fn); v {
		t.Error("expected false on error")
	}
}

func TestPollTickReturnsCmd(t *testing.T) {
	cmd := PollTick(time.Second)
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}
}

func TestRouterSettingsZeroValue(t *testing.T) {
	s := RouterSettings{}
	if s.BWIn != "" || s.BWOut != "" || s.BWShare != "" {
		t.Error("expected zero-value strings")
	}
}

func TestTokenMutexReadWrite(t *testing.T) {
	tokenMu.Lock()
	authToken = "test-token"
	tokenMu.Unlock()

	tokenMu.RLock()
	tok := authToken
	tokenMu.RUnlock()

	if tok != "test-token" {
		t.Errorf("expected test-token, got %s", tok)
	}

	// Reset for other tests.
	tokenMu.Lock()
	authToken = ""
	tokenMu.Unlock()
}
