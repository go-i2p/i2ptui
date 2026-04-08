package i2ptui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Padding(0, 2)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	labelStyle = lipgloss.NewStyle().
			Width(20).
			Foreground(lipgloss.Color("252"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	confirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 3)

	notifyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("228")).
			Bold(true)
)

// fmtDuration formats a duration as a human-readable string (e.g. "2d 3h 15m").
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// fmtBandwidth formats bytes-per-second as a human-readable string.
func fmtBandwidth(bps int) string {
	if bps >= 1024*1024 {
		return fmt.Sprintf("%.1f MB/s", float64(bps)/(1024*1024))
	}
	if bps >= 1024 {
		return fmt.Sprintf("%.1f KB/s", float64(bps)/1024)
	}
	return fmt.Sprintf("%d B/s", bps)
}
