// Command i2ptui is the standalone CLI for the I2P router TUI.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-i2p/i2ptui"
)

func main() {
	// Load persistent config as defaults.
	cfg := i2ptui.LoadConfig(i2ptui.DefaultConfigPath())

	host := flag.String("host", defaultStr(cfg.Host, "127.0.0.1"), "I2PControl host")
	port := flag.String("port", defaultStr(cfg.Port, "7657"), "I2PControl port")
	path := flag.String("path", defaultStr(cfg.Path, "jsonrpc"), "RPC URL path")
	password := flag.String("password", defaultStr(cfg.Password, "itoopie"), "I2PControl API password")
	cert := flag.String("cert", cfg.Cert, "path to self-signed cert")
	interval := flag.Duration("interval", 5*time.Second, "polling interval")
	theme := flag.String("theme", defaultStr(cfg.Theme, "dark"), "color theme (dark, light)")
	flag.Parse()

	var opts []i2ptui.Option
	opts = append(opts,
		i2ptui.WithHost(*host),
		i2ptui.WithPort(*port),
		i2ptui.WithPath(*path),
		i2ptui.WithPassword(*password),
		i2ptui.WithCert(*cert),
		i2ptui.WithInterval(*interval),
	)
	if *theme == "light" {
		opts = append(opts, i2ptui.WithTheme(i2ptui.LightTheme))
	}

	m := i2ptui.New(opts...)

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// defaultStr returns val if non-empty, otherwise fallback.
func defaultStr(val, fallback string) string {
	if val != "" {
		return val
	}
	return fallback
}
