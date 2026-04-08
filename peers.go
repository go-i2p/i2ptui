package i2ptui

import (
	"fmt"
	"strings"

	"github.com/go-i2p/i2ptui/rpc"
)

type peersModel struct {
	snapshot rpc.RouterSnapshot
}

// newPeersModel returns an empty peersModel.
func newPeersModel() peersModel {
	return peersModel{}
}

// SetSnapshot updates the snapshot data used by the peers table.
func (m peersModel) SetSnapshot(s rpc.RouterSnapshot) peersModel {
	m.snapshot = s
	return m
}

// View renders the peers tab with tunnel and peer counts.
func (m peersModel) View(width int) string {
	s := m.snapshot
	var b strings.Builder

	b.WriteString(sectionStyle.Render("  Peer Overview"))
	b.WriteString("\n")
	b.WriteString(m.row("Known Peers", fmt.Sprintf("%d", s.KnownPeers)))
	b.WriteString(m.row("Reseeding", fmt.Sprintf("%v", s.Reseeding)))
	b.WriteString(m.row("Net Status", s.NetStatus))

	return b.String()
}

// row formats a label-value pair for display.
func (m peersModel) row(label, value string) string {
	return fmt.Sprintf("  %s %s\n",
		labelStyle.Render(label),
		valueStyle.Render(value),
	)
}
