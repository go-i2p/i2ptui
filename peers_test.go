package i2ptui

import (
	"strings"
	"testing"

	"github.com/go-i2p/i2ptui/rpc"
)

func TestPeersView(t *testing.T) {
	m := newPeersModel()
	m = m.SetSnapshot(rpc.RouterSnapshot{
		KnownPeers: 1842,
		Reseeding:  false,
		NetStatus:  "OK",
	})

	v := m.View(80)

	checks := []string{"1842", "false", "OK", "Peer Overview"}
	for _, c := range checks {
		if !strings.Contains(v, c) {
			t.Errorf("peers view missing %q", c)
		}
	}
}

func TestPeersViewEmpty(t *testing.T) {
	m := newPeersModel()
	v := m.View(80)
	if !strings.Contains(v, "Peer Overview") {
		t.Error("expected Peer Overview section header")
	}
}
