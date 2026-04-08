package i2ptui

import (
	"sync"
	"time"
)

// historyBuffer stores time-series samples within a rolling window.
type historyBuffer struct {
	mu      sync.Mutex
	samples []sample
	window  time.Duration
}

type sample struct {
	time  time.Time
	value float64
}

// newHistoryBuffer creates a buffer that retains samples within the given window.
func newHistoryBuffer(window time.Duration) *historyBuffer {
	return &historyBuffer{
		window: window,
	}
}

// Add appends a sample and discards entries older than the window.
func (h *historyBuffer) Add(t time.Time, v float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.samples = append(h.samples, sample{time: t, value: v})
	h.expire(t)
}

// Values returns the current sample values in chronological order.
func (h *historyBuffer) Values() []float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.expire(time.Now())
	out := make([]float64, len(h.samples))
	for i, s := range h.samples {
		out[i] = s.value
	}
	return out
}

// Len returns the number of current samples.
func (h *historyBuffer) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.samples)
}

// expire removes samples older than the window relative to ref.
func (h *historyBuffer) expire(ref time.Time) {
	cutoff := ref.Add(-h.window)
	i := 0
	for i < len(h.samples) && h.samples[i].time.Before(cutoff) {
		i++
	}
	if i > 0 {
		h.samples = h.samples[i:]
	}
}
