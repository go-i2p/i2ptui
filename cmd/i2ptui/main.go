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
	host := flag.String("host", "127.0.0.1", "I2PControl host")
	port := flag.String("port", "7657", "I2PControl port")
	path := flag.String("path", "jsonrpc", "RPC URL path")
	password := flag.String("password", "itoopie", "I2PControl API password")
	cert := flag.String("cert", "", "path to self-signed cert")
	interval := flag.Duration("interval", 5*time.Second, "polling interval")
	flag.Parse()

	m := i2ptui.New(
		i2ptui.WithHost(*host),
		i2ptui.WithPort(*port),
		i2ptui.WithPath(*path),
		i2ptui.WithPassword(*password),
		i2ptui.WithCert(*cert),
		i2ptui.WithInterval(*interval),
	)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
