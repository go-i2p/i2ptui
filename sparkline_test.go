package i2ptui

import (
	"strings"
	"testing"
)

func TestRenderSparklineEmpty(t *testing.T) {
	if got := renderSparkline(nil, 0); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestRenderSparklineSingle(t *testing.T) {
	got := renderSparkline([]float64{42}, 0)
	if len([]rune(got)) != 1 {
		t.Errorf("expected 1 rune, got %d", len([]rune(got)))
	}
}

func TestRenderSparklineAscending(t *testing.T) {
	vals := []float64{0, 1, 2, 3, 4, 5, 6, 7}
	got := renderSparkline(vals, 0)
	runes := []rune(got)
	if len(runes) != 8 {
		t.Fatalf("expected 8 runes, got %d", len(runes))
	}
	if runes[0] != '▁' {
		t.Errorf("first rune should be lowest block, got %c", runes[0])
	}
	if runes[7] != '█' {
		t.Errorf("last rune should be highest block, got %c", runes[7])
	}
}

func TestRenderSparklineWidth(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	got := renderSparkline(vals, 5)
	runes := []rune(got)
	if len(runes) != 5 {
		t.Errorf("expected 5 runes with width=5, got %d", len(runes))
	}
}

func TestRenderSparklineConstant(t *testing.T) {
	vals := []float64{5, 5, 5, 5}
	got := renderSparkline(vals, 0)
	runes := []rune(got)
	// All same value → all same block
	for i := 1; i < len(runes); i++ {
		if runes[i] != runes[0] {
			t.Errorf("constant values should all map to same block")
			break
		}
	}
}

func TestRenderBarChart(t *testing.T) {
	entries := []barEntry{
		{label: "A", value: 100},
		{label: "B", value: 50},
		{label: "C", value: 0},
	}
	got := renderBarChart(entries, 10)
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A ") {
		t.Errorf("first line should start with 'A ', got %q", lines[0])
	}
	if !strings.Contains(lines[0], "██████████") {
		t.Errorf("expected full bar for max value, got %q", lines[0])
	}
}

func TestRenderBarChartEmpty(t *testing.T) {
	if got := renderBarChart(nil, 10); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
