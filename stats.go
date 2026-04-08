package i2ptui

import (
	"fmt"
	"strings"

	"github.com/go-i2p/i2ptui/rpc"
)

type statsModel struct {
	snapshot rpc.RouterSnapshot
}

// newStatsModel returns an empty statsModel.
func newStatsModel() statsModel {
	return statsModel{}
}

// SetSnapshot updates the snapshot data.
func (m statsModel) SetSnapshot(s rpc.RouterSnapshot) statsModel {
	m.snapshot = s
	return m
}

// View renders the stats tab.
func (m statsModel) View(width int) string {
	s := m.snapshot
	var b strings.Builder

	b.WriteString(sectionStyle.Render("  Bandwidth Stats"))
	b.WriteString("\n")
	b.WriteString(m.row("Send (bps)", fmt.Sprintf("%d", s.SendBps)))
	b.WriteString(m.row("Receive (bps)", fmt.Sprintf("%d", s.ReceiveBps)))
	b.WriteString(m.row("Send Avg/hr", fmt.Sprintf("%d", s.SendBpsHourAvg)))
	b.WriteString(m.row("Receive Avg/hr", fmt.Sprintf("%d", s.ReceiveBpsHourAvg)))

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("  Tunnel Build Stats"))
	b.WriteString("\n")
	b.WriteString(m.row("Exploratory Success", fmt.Sprintf("%d", s.ExplBuildSuccess)))
	b.WriteString(m.row("Exploratory Reject", fmt.Sprintf("%d", s.ExplBuildReject)))
	b.WriteString(m.row("Exploratory Expire", fmt.Sprintf("%d", s.ExplBuildExpire)))
	b.WriteString(m.row("Client Build Success", fmt.Sprintf("%d", s.ClientBuildSuccess)))
	b.WriteString(m.row("Build Request Time", fmt.Sprintf("%.2fs", s.BuildRequestTime)))

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("  Tunnel Build Percentages"))
	b.WriteString("\n")
	b.WriteString(m.row("Success %", fmt.Sprintf("%d%%", s.ExplBuildSuccessPct)))
	b.WriteString(m.row("Reject %", fmt.Sprintf("%d%%", s.ExplBuildRejectPct)))
	b.WriteString(m.row("Expire %", fmt.Sprintf("%d%%", s.ExplBuildExpirePct)))

	return b.String()
}

// row formats a label-value pair for display.
func (m statsModel) row(label, value string) string {
	return fmt.Sprintf("  %s %s\n",
		labelStyle.Render(label),
		valueStyle.Render(value),
	)
}
