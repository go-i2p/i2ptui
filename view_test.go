package i2ptui

import (
	"strings"
	"testing"

	"github.com/go-i2p/i2ptui/rpc"
)

func TestOverviewView(t *testing.T) {
	m := newOverviewModel()
	m = m.SetSnapshot(rpc.RouterSnapshot{
		Status:               "OK",
		NetStatus:            "OK",
		Version:              "0.9.62",
		Uptime:               3600000,
		IncomingBW:           51200,
		OutgoingBW:           12800,
		KnownPeers:           1842,
		ParticipatingTunnels: 352,
	})

	v := m.View(80)

	checks := []string{"OK", "0.9.62", "1842", "352"}
	for _, c := range checks {
		if !strings.Contains(v, c) {
			t.Errorf("overview view missing %q", c)
		}
	}
}

func TestStatsView(t *testing.T) {
	m := newStatsModel()
	m = m.SetSnapshot(rpc.RouterSnapshot{
		SendBps:             12400,
		ReceiveBps:          48200,
		ExplBuildSuccess:    72,
		ExplBuildReject:     18,
		ExplBuildExpire:     10,
		ExplBuildSuccessPct: 72,
		ExplBuildRejectPct:  18,
		ExplBuildExpirePct:  10,
		ClientBuildSuccess:  54,
		BuildRequestTime:    1.23,
	})

	v := m.View(80)

	checks := []string{"12400", "48200", "72%", "18%", "10%", "1.23"}
	for _, c := range checks {
		if !strings.Contains(v, c) {
			t.Errorf("stats view missing %q", c)
		}
	}
}
