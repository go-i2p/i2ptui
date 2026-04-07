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
	tabControl
)

var tabNames = []string{"Dashboard", "Stats", "Control"}

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
	control  controlModel

	statusMsg  string
	statusTime time.Time
}

// New creates a new Model with the given options.
func New(opts ...Option) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := Model{
		host:     "127.0.0.1",
		port:     "7657",
		path:     "jsonrpc",
		password: "itoopie",
		interval: 5 * time.Second,
		spinner:  s,
		loading:  true,
		overview: newOverviewModel(),
		stats:    newStatsModel(),
		control:  newControlModel(),
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
			m.activeTab = tabControl
			return m, nil
		case "r":
			m.loading = true
			return m, rpc.FetchSnapshotCmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case rpc.RouterSnapshot:
		m.snapshot = msg
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.err = nil
		}
		m.overview = m.overview.SetSnapshot(msg)
		m.stats = m.stats.SetSnapshot(msg)
		return m, rpc.PollTick(m.interval)

	case rpc.TickMsg:
		return m, rpc.FetchSnapshotCmd

	case controlActionMsg:
		m.statusMsg = msg.result
		m.statusTime = time.Now()
		return m, nil

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
	case tabStats:
		b.WriteString(m.stats.View(m.width))
	case tabControl:
		b.WriteString(m.control.View(m.width))
	}

	if m.statusMsg != "" && time.Since(m.statusTime) < 10*time.Second {
		b.WriteString("\n")
		b.WriteString(statusStyle.Render(fmt.Sprintf("  %s", m.statusMsg)))
	}

	b.WriteString(m.renderFooter())

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
	help := "tab: switch  r: refresh  q: quit"
	gap := m.width - lipgloss.Width(age) - lipgloss.Width(help)
	if gap < 1 {
		gap = 1
	}
	return "\n" + footerStyle.Render(help+strings.Repeat(" ", gap)+age)
}
