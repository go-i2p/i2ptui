package i2ptui

import "math"

// sparkline characters from lowest to highest.
var sparkBlocks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// renderSparkline returns a Unicode sparkline string for the given values.
// width limits the number of trailing samples used; 0 means use all.
func renderSparkline(values []float64, width int) string {
	if len(values) == 0 {
		return ""
	}
	if width > 0 && len(values) > width {
		values = values[len(values)-width:]
	}

	min, max := minMax(values)
	span := max - min
	if span == 0 {
		span = 1
	}

	out := make([]rune, len(values))
	for i, v := range values {
		idx := int(math.Round((v - min) / span * float64(len(sparkBlocks)-1)))
		out[i] = sparkBlocks[clampIndex(idx, len(sparkBlocks))]
	}
	return string(out)
}

// renderBarChart produces a simple horizontal bar chart.
// Each entry is label + bar of width proportional to value.
func renderBarChart(entries []barEntry, maxWidth int) string {
	if len(entries) == 0 || maxWidth <= 0 {
		return ""
	}
	var maxVal float64
	for _, e := range entries {
		if e.value > maxVal {
			maxVal = e.value
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	var result string
	for _, e := range entries {
		barLen := int(math.Round(e.value / maxVal * float64(maxWidth)))
		if barLen < 0 {
			barLen = 0
		}
		bar := make([]rune, barLen)
		for i := range bar {
			bar[i] = '█'
		}
		result += e.label + " " + string(bar) + "\n"
	}
	return result
}

type barEntry struct {
	label string
	value float64
}

// minMax returns the minimum and maximum values in vals.
func minMax(vals []float64) (float64, float64) {
	mn, mx := vals[0], vals[0]
	for _, v := range vals[1:] {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return mn, mx
}

// clampIndex constrains idx to the range [0, length-1].
func clampIndex(idx, length int) int {
	if idx < 0 {
		return 0
	}
	if idx >= length {
		return length - 1
	}
	return idx
}
