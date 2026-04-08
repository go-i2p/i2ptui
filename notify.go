package i2ptui

import (
	"fmt"
	"time"

	"github.com/go-i2p/i2ptui/rpc"
)

// notification represents a dismissable in-TUI notification.
type notification struct {
	message string
	time    time.Time
}

type notifyModel struct {
	notifications []notification
	prevStatus    string
	prevNetStatus string
	prevReseeding bool
}

// newNotifyModel returns a notifyModel with no initial notifications.
func newNotifyModel() notifyModel {
	return notifyModel{}
}

// CheckChanges compares old and new snapshots and adds notifications.
func (m *notifyModel) CheckChanges(prev, curr rpc.RouterSnapshot) {
	if m.prevStatus != "" && m.prevStatus != curr.Status {
		m.add(fmt.Sprintf("Status changed: %s → %s", m.prevStatus, curr.Status))
	}
	if m.prevNetStatus != "" && m.prevNetStatus != curr.NetStatus {
		m.add(fmt.Sprintf("Net status changed: %s → %s", m.prevNetStatus, curr.NetStatus))
	}
	if m.prevStatus != "" && !m.prevReseeding && curr.Reseeding {
		m.add("Reseed triggered")
	}
	m.prevStatus = curr.Status
	m.prevNetStatus = curr.NetStatus
	m.prevReseeding = curr.Reseeding
}

// add appends a notification and sends a desktop notification.
func (m *notifyModel) add(msg string) {
	m.notifications = append(m.notifications, notification{
		message: msg,
		time:    time.Now(),
	})
	desktopNotify("i2ptui", msg)
}

// Dismiss removes the first notification.
func (m *notifyModel) Dismiss() {
	if len(m.notifications) > 0 {
		m.notifications = m.notifications[1:]
	}
}

// HasNotifications returns true if there are pending notifications.
func (m *notifyModel) HasNotifications() bool {
	return len(m.notifications) > 0
}

// View renders the notification bar.
func (m *notifyModel) View() string {
	if len(m.notifications) == 0 {
		return ""
	}
	n := m.notifications[0]
	age := fmtDuration(time.Since(n.time))
	return notifyStyle.Render(fmt.Sprintf(
		"  ● %s (%s ago)  [Esc to dismiss]", n.message, age,
	))
}
