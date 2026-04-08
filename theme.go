package i2ptui

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the TUI.
type Theme struct {
	ActiveTab   lipgloss.Color
	InactiveTab lipgloss.Color
	Header      lipgloss.Color
	Label       lipgloss.Color
	Value       lipgloss.Color
	Section     lipgloss.Color
	Error       lipgloss.Color
	Status      lipgloss.Color
	Footer      lipgloss.Color
	Selected    lipgloss.Color
	ConfirmBdr  lipgloss.Color
	Notify      lipgloss.Color
	Restart     lipgloss.Color
}

// DarkTheme is the default dark color scheme.
var DarkTheme = Theme{
	ActiveTab:   lipgloss.Color("205"),
	InactiveTab: lipgloss.Color("240"),
	Header:      lipgloss.Color("240"),
	Label:       lipgloss.Color("252"),
	Value:       lipgloss.Color("86"),
	Section:     lipgloss.Color("99"),
	Error:       lipgloss.Color("196"),
	Status:      lipgloss.Color("220"),
	Footer:      lipgloss.Color("240"),
	Selected:    lipgloss.Color("205"),
	ConfirmBdr:  lipgloss.Color("205"),
	Notify:      lipgloss.Color("228"),
	Restart:     lipgloss.Color("214"),
}

// LightTheme is a lighter color scheme for light terminals.
var LightTheme = Theme{
	ActiveTab:   lipgloss.Color("127"),
	InactiveTab: lipgloss.Color("245"),
	Header:      lipgloss.Color("245"),
	Label:       lipgloss.Color("238"),
	Value:       lipgloss.Color("28"),
	Section:     lipgloss.Color("63"),
	Error:       lipgloss.Color("160"),
	Status:      lipgloss.Color("172"),
	Footer:      lipgloss.Color("245"),
	Selected:    lipgloss.Color("127"),
	ConfirmBdr:  lipgloss.Color("127"),
	Notify:      lipgloss.Color("136"),
	Restart:     lipgloss.Color("166"),
}

// applyTheme updates the global styles with the given theme's colors.
func applyTheme(t Theme) {
	activeTabStyle = activeTabStyle.Foreground(t.ActiveTab)
	inactiveTabStyle = inactiveTabStyle.Foreground(t.InactiveTab)
	headerStyle = headerStyle.Foreground(t.Header)
	labelStyle = labelStyle.Foreground(t.Label)
	valueStyle = valueStyle.Foreground(t.Value)
	sectionStyle = sectionStyle.Foreground(t.Section)
	errorStyle = errorStyle.Foreground(t.Error)
	statusStyle = statusStyle.Foreground(t.Status)
	footerStyle = footerStyle.Foreground(t.Footer)
	selectedStyle = selectedStyle.Foreground(t.Selected)
	confirmBoxStyle = confirmBoxStyle.BorderForeground(t.ConfirmBdr)
	notifyStyle = notifyStyle.Foreground(t.Notify)
	restartStyle = restartStyle.Foreground(t.Restart)
}
