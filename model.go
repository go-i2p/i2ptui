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

// WithTheme sets the color theme.
func WithTheme(t Theme) Option {
	return func(m *Model) { applyTheme(t) }
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

// initRPC connects to I2PControl and fetches the first snapshot.
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
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case rpc.RouterSnapshot:
		return m.handleSnapshot(msg)
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
	return m.delegateToTab(msg)
}

// numKeyTab maps numeric key strings to tab indices.
var numKeyTab = map[string]int{
	"1": tabDashboard,
	"2": tabStats,
	"3": tabPeers,
	"4": tabControl,
	"5": tabSettings,
}

// handleKey processes keyboard input for the root model.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	key := msg.String()
	if tab, ok := numKeyTab[key]; ok {
		m.activeTab = tab
		if tab == tabSettings && !m.settings.loaded {
			return m, fetchSettingsCmd
		}
		return m, nil
	}
	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab":
		m.activeTab = (m.activeTab + 1) % len(tabNames)
	case "shift+tab":
		m.activeTab = (m.activeTab - 1 + len(tabNames)) % len(tabNames)
	case "r":
		m.loading = true
		return m, rpc.FetchSnapshotCmd
	case "g":
		m.showGraphs = !m.showGraphs
	case "esc":
		if m.notify.HasNotifications() {
			m.notify.Dismiss()
		}
	}
	return m, nil
}

// handleMouse processes mouse input for tab switching.
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
		if msg.Y == 0 {
			m.activeTab = m.tabFromClick(msg.X)
		}
	}
	return m, nil
}

// handleSnapshot applies a new router snapshot to the model.
func (m Model) handleSnapshot(msg rpc.RouterSnapshot) (tea.Model, tea.Cmd) {
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
}

// delegateToTab forwards unhandled messages to the active tab's sub-model.
func (m Model) delegateToTab(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	b.WriteString(m.renderActiveTab())
	b.WriteString(m.renderStatusBar())
	b.WriteString(m.renderFooter())

	return b.String()
}

// renderActiveTab returns the view string for the currently selected tab.
func (m Model) renderActiveTab() string {
	switch m.activeTab {
	case tabDashboard:
		s := m.overview.View(m.width)
		if m.showGraphs {
			s += m.renderGraphs()
		}
		return s
	case tabStats:
		s := m.stats.View(m.width)
		if m.showGraphs {
			s += m.renderBuildChart()
		}
		return s
	case tabPeers:
		s := m.peers.View(m.width)
		if m.showGraphs {
			s += m.renderPeerGraph()
		}
		return s
	case tabControl:
		return m.control.View(m.width)
	case tabSettings:
		return m.settings.View(m.width)
	default:
		return ""
	}
}

// renderStatusBar returns the status message and notification bar.
func (m Model) renderStatusBar() string {
	var b strings.Builder
	if m.statusMsg != "" && time.Since(m.statusTime) < 10*time.Second {
		b.WriteString("\n")
		b.WriteString(statusStyle.Render(fmt.Sprintf("  %s", m.statusMsg)))
	}
	if nv := m.notify.View(); nv != "" {
		b.WriteString("\n")
		b.WriteString(nv)
	}
	return b.String()
}

// recordHistory appends the snapshot's metrics to their history buffers.
func (m *Model) recordHistory(snap rpc.RouterSnapshot) {
	t := snap.FetchedAt
	m.inBWHistory.Add(t, float64(snap.IncomingBW))
	m.outBWHistory.Add(t, float64(snap.OutgoingBW))
	m.tunnelHistory.Add(t, float64(snap.ParticipatingTunnels))
	m.buildSuccHistory.Add(t, float64(snap.ExplBuildSuccessPct))
	m.peerHistory.Add(t, float64(snap.KnownPeers))
}

// renderGraphs returns the sparkline graph panel for the Dashboard tab.
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

// graphRow formats a single sparkline row with a label.
func graphRow(label string, vals []float64, width int) string {
	spark := renderSparkline(vals, width)
	if spark == "" {
		spark = "(no data)"
	}
	return fmt.Sprintf("  %s %s\n", labelStyle.Render(label), spark)
}

// renderBuildChart returns the build success bar chart for the Stats tab.
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

// renderPeerGraph returns the peer count sparkline for the Peers tab.
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

// renderTabs returns the tab bar header.
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

// renderFooter returns the bottom help and data-age line.
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

// tabFromClick returns the tab index for a mouse click at column x.
func (m Model) tabFromClick(x int) int {
	offset := 0
	for i, name := range tabNames {
		w := lipgloss.Width(name) + 4 // padding 0,2 on each side
		if x >= offset && x < offset+w {
			return i
		}
		offset += w
	}
	return m.activeTab
}
