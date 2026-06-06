package i2ptui

import (
	"fmt"
	"strings"

	"github.com/go-i2p/i2ptui/rpc"
)

type overviewModel struct {
	snapshot rpc.RouterSnapshot
}

// newOverviewModel returns an empty overviewModel.
func newOverviewModel() overviewModel {
	return overviewModel{}
}

// SetSnapshot updates the snapshot data used by the dashboard.
func (m overviewModel) SetSnapshot(s rpc.RouterSnapshot) overviewModel {
	m.snapshot = s
	return m
}

// View renders the dashboard tab with the router overview.
func (m overviewModel) View(width int) string {
	s := m.snapshot
	var b strings.Builder

	b.WriteString(sectionStyle.Render("  Router Info"))
	b.WriteString("\n")
	b.WriteString(m.row("Status", s.Status))
	b.WriteString(m.row("Net Status", s.NetStatus))
	b.WriteString(m.row("Version", s.Version))
	b.WriteString(m.row("Router Hash", s.RouterHash))
	b.WriteString(m.row("Uptime", fmtDuration(s.UptimeDuration())))
	b.WriteString(m.row("Known Peers", fmt.Sprintf("%d", s.KnownPeers)))
	b.WriteString(m.row("Reseeding", fmt.Sprintf("%v", s.Reseeding)))

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("  Bandwidth"))
	b.WriteString("\n")
	b.WriteString(m.row("Incoming", fmtBandwidth(s.IncomingBW)))
	b.WriteString(m.row("Outgoing", fmtBandwidth(s.OutgoingBW)))
	b.WriteString(m.row("Receive Avg/hr", fmtBandwidth(s.ReceiveBpsHourAvg)))
	b.WriteString(m.row("Send Avg/hr", fmtBandwidth(s.SendBpsHourAvg)))

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("  Tunnels"))
	b.WriteString("\n")
	b.WriteString(m.row("Participating", fmt.Sprintf("%d", s.ParticipatingTunnels)))
	b.WriteString(m.row("Avg (5 min)", fmt.Sprintf("%d", s.ParticipatingAvg)))
	b.WriteString(m.row("Avg (1 hour)", fmt.Sprintf("%d", s.ParticipatingHourAvg)))

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("  Settings"))
	b.WriteString("\n")
	b.WriteString(m.row("UPnP", s.Upnp))

	return b.String()
}

// row formats a label-value pair for display.
func (m overviewModel) row(label, value string) string {
	return fmt.Sprintf("  %s %s\n",
		labelStyle.Render(label),
		valueStyle.Render(value),
	)
}
