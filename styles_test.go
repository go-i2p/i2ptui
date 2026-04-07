package i2ptui

import (
	"testing"
	"time"
)

func TestFmtDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{5 * time.Second, "5s"},
		{90 * time.Second, "1m 30s"},
		{3661 * time.Second, "1h 1m 1s"},
		{90061 * time.Second, "1d 1h 1m"},
	}
	for _, tt := range tests {
		got := fmtDuration(tt.d)
		if got != tt.want {
			t.Errorf("fmtDuration(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestFmtBandwidth(t *testing.T) {
	tests := []struct {
		bps  int
		want string
	}{
		{500, "500 B/s"},
		{2048, "2.0 KB/s"},
		{1572864, "1.5 MB/s"},
	}
	for _, tt := range tests {
		got := fmtBandwidth(tt.bps)
		if got != tt.want {
			t.Errorf("fmtBandwidth(%d) = %q, want %q", tt.bps, got, tt.want)
		}
	}
}
