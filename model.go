package i2ptui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-i2p/i2ptui/rpc"
)

// Tab identifiers.
const (
	tabDashboard = iota
	tabStats
	tabPeers
	tabControl
	tabSettings
)

var tabNames = []string{"Dashboard", "Stats", "Peers", "Control", "Settings"}

// Option configures a Model.
type Option func(*Model)

// WithHost sets the I2PControl host.
func WithHost(h string) Option {
	return func(m *Model) { m.host = h }
}

// WithPort sets the I2PControl port.
func WithPort(p string) Option {
	return func(m *Model) { m.port = p }
}

// WithPath sets the RPC URL path.
func WithPath(p string) Option {
	return func(m *Model) { m.path = p }
}

// WithPassword sets the I2PControl API password.
func WithPassword(p string) Option {
	return func(m *Model) { m.password = p }
}

// WithCert sets the path to a self-signed certificate.
func WithCert(c string) Option {
	return func(m *Model) { m.cert = c }
}

// WithInterval sets the polling interval.
func WithInterval(d time.Duration) Option {
	return func(m *Model) { m.interval = d }
}

// Model is the root Bubble Tea model for i2ptui.
type Model struct {
	host     string
	port     string
	path     string
	password string
	cert     string
	interval time.Duration

	activeTab int
	width     int
	height    int

	snapshot rpc.RouterSnapshot
	spinner  spinner.Model
	loading  bool
	err      error

	overview overviewModel
	stats    statsModel
	peers    peersModel
	control  controlModel
	settings settingsModel

	statusMsg  string
	statusTime time.Time

	notify notifyModel

	showGraphs       bool
	inBWHistory      *historyBuffer
	outBWHistory     *historyBuffer
	tunnelHistory    *historyBuffer
	buildSuccHistory *historyBuffer
	peerHistory      *historyBuffer
}

// New creates a new Model with the given options.
func New(opts ...Option) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	histWindow := 5 * time.Minute
	m := Model{
		host:             "127.0.0.1",
		port:             "7657",
		path:             "jsonrpc",
		password:         "itoopie",
		interval:         5 * time.Second,
		spinner:          s,
		loading:          true,
		overview:         newOverviewModel(),
		stats:            newStatsModel(),
		peers:            newPeersModel(),
		control:          newControlModel(),
		settings:         newSettingsModel(),
		notify:           newNotifyModel(),
		inBWHistory:      newHistoryBuffer(histWindow),
		outBWHistory:     newHistoryBuffer(histWindow),
		tunnelHistory:    newHistoryBuffer(histWindow),
		buildSuccHistory: newHistoryBuffer(histWindow),
		peerHistory:      newHistoryBuffer(histWindow),
	}
	for _, o := range opts {
		o(&m)
	}
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.initRPC,
	)
}

func (m Model) initRPC() tea.Msg {
	err := rpc.Setup(m.host, m.port, m.path, m.password, m.cert)
	if err != nil {
		return rpc.RouterSnapshot{Err: err, FetchedAt: time.Now()}
	}
	return rpc.FetchSnapshot()
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.control.confirming {
			var cmd tea.Cmd
			m.control, cmd = m.control.Update(msg)
			return m, cmd
		}
		if m.settings.confirming {
			var cmd tea.Cmd
			m.settings, cmd = m.settings.Update(msg)
			return m, cmd
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.activeTab = (m.activeTab + 1) % len(tabNames)
			return m, nil
		case "shift+tab":
			m.activeTab = (m.activeTab - 1 + len(tabNames)) % len(tabNames)
			return m, nil
		case "1":
			m.activeTab = tabDashboard
			return m, nil
		case "2":
			m.activeTab = tabStats
			return m, nil
		case "3":
			m.activeTab = tabPeers
			return m, nil
		case "4":
			m.activeTab = tabControl
			return m, nil
		case "5":
			m.activeTab = tabSettings
			if !m.settings.loaded {
				return m, fetchSettingsCmd
			}
			return m, nil
		case "r":
			m.loading = true
			return m, rpc.FetchSnapshotCmd
		case "g":
			m.showGraphs = !m.showGraphs
			return m, nil
		case "esc":
			if m.notify.HasNotifications() {
				m.notify.Dismiss()
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case rpc.RouterSnapshot:
		prevSnap := m.snapshot
		m.snapshot = msg
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.err = nil
			m.recordHistory(msg)
			m.notify.CheckChanges(prevSnap, msg)
		}
		m.overview = m.overview.SetSnapshot(msg)
		m.stats = m.stats.SetSnapshot(msg)
		m.peers = m.peers.SetSnapshot(msg)
		return m, rpc.PollTick(m.interval)

	case rpc.TickMsg:
		return m, rpc.FetchSnapshotCmd

	case controlActionMsg:
		m.statusMsg = msg.result
		m.statusTime = time.Now()
		return m, nil

	case settingsMsg:
		var cmd tea.Cmd
		m.settings, cmd = m.settings.Update(msg)
		return m, cmd

	case settingsSavedMsg:
		var cmd tea.Cmd
		m.settings, cmd = m.settings.Update(msg)
		return m, cmd

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.activeTab == tabControl {
		var cmd tea.Cmd
		m.control, cmd = m.control.Update(msg)
		return m, cmd
	}

	if m.activeTab == tabSettings {
		var cmd tea.Cmd
		m.settings, cmd = m.settings.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var b strings.Builder

	b.WriteString(m.renderTabs())
	b.WriteString("\n")
	b.WriteString(headerStyle.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")

	if m.loading && m.snapshot.FetchedAt.IsZero() {
		b.WriteString(fmt.Sprintf("\n  %s Connecting to I2PControl...\n", m.spinner.View()))
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("\n  Error: %v\n", m.err)))
		b.WriteString("\n  Press 'r' to retry, 'q' to quit.\n")
		return b.String()
	}

	switch m.activeTab {
	case tabDashboard:
		b.WriteString(m.overview.View(m.width))
		if m.showGraphs {
			b.WriteString(m.renderGraphs())
		}
	case tabStats:
		b.WriteString(m.stats.View(m.width))
		if m.showGraphs {
			b.WriteString(m.renderBuildChart())
		}
	case tabPeers:
		b.WriteString(m.peers.View(m.width))
		if m.showGraphs {
			b.WriteString(m.renderPeerGraph())
		}
	case tabControl:
		b.WriteString(m.control.View(m.width))
	case tabSettings:
		b.WriteString(m.settings.View(m.width))
	}

	if m.statusMsg != "" && time.Since(m.statusTime) < 10*time.Second {
		b.WriteString("\n")
		b.WriteString(statusStyle.Render(fmt.Sprintf("  %s", m.statusMsg)))
	}

	if nv := m.notify.View(); nv != "" {
		b.WriteString("\n")
		b.WriteString(nv)
	}

	b.WriteString(m.renderFooter())

	return b.String()
}

func (m *Model) recordHistory(snap rpc.RouterSnapshot) {
	t := snap.FetchedAt
	m.inBWHistory.Add(t, float64(snap.IncomingBW))
	m.outBWHistory.Add(t, float64(snap.OutgoingBW))
	m.tunnelHistory.Add(t, float64(snap.ParticipatingTunnels))
	m.buildSuccHistory.Add(t, float64(snap.ExplBuildSuccessPct))
	m.peerHistory.Add(t, float64(snap.KnownPeers))
}

func (m Model) renderGraphs() string {
	var b strings.Builder
	sparkWidth := m.width - 25
	if sparkWidth < 10 {
		sparkWidth = 10
	}
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("  Graphs"))
	b.WriteString("\n")
	b.WriteString(graphRow("Incoming BW", m.inBWHistory.Values(), sparkWidth))
	b.WriteString(graphRow("Outgoing BW", m.outBWHistory.Values(), sparkWidth))
	b.WriteString(graphRow("Tunnels", m.tunnelHistory.Values(), sparkWidth))
	return b.String()
}

func graphRow(label string, vals []float64, width int) string {
	spark := renderSparkline(vals, width)
	if spark == "" {
		spark = "(no data)"
	}
	return fmt.Sprintf("  %s %s\n", labelStyle.Render(label), spark)
}

func (m Model) renderBuildChart() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("  Build Success Rate"))
	b.WriteString("\n")
	s := m.snapshot
	entries := []barEntry{
		{label: "Success", value: float64(s.ExplBuildSuccessPct)},
		{label: "Reject ", value: float64(s.ExplBuildRejectPct)},
		{label: "Expire ", value: float64(s.ExplBuildExpirePct)},
	}
	chartWidth := m.width - 20
	if chartWidth < 10 {
		chartWidth = 10
	}
	b.WriteString(renderBarChart(entries, chartWidth))
	return b.String()
}

func (m Model) renderPeerGraph() string {
	var b strings.Builder
	sparkWidth := m.width - 25
	if sparkWidth < 10 {
		sparkWidth = 10
	}
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("  Peer Count History"))
	b.WriteString("\n")
	b.WriteString(graphRow("Known Peers", m.peerHistory.Values(), sparkWidth))
	return b.String()
}

func (m Model) renderTabs() string {
	var tabs []string
	for i, name := range tabNames {
		if i == m.activeTab {
			tabs = append(tabs, activeTabStyle.Render(name))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderFooter() string {
	age := ""
	if !m.snapshot.FetchedAt.IsZero() {
		age = fmt.Sprintf("Updated %s ago", fmtDuration(time.Since(m.snapshot.FetchedAt)))
	}
	help := "tab: switch  r: refresh  g: graphs  q: quit"
	gap := m.width - lipgloss.Width(age) - lipgloss.Width(help)
	if gap < 1 {
		gap = 1
	}
	return "\n" + footerStyle.Render(help+strings.Repeat(" ", gap)+age)
}
