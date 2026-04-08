// Command i2ptui is the standalone CLI for the I2P router TUI.
package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-i2p/i2ptui"
	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags.
var version = "dev"

// rootCmd is the top-level cobra command for i2ptui.
var rootCmd = &cobra.Command{
	Use:   "i2ptui",
	Short: "Terminal UI for monitoring and controlling an I2P router",
	Long: `i2ptui is a terminal user interface for the I2P router.
It connects via I2PControl JSON-RPC to display router status,
bandwidth stats, peer information, and router settings.`,
	SilenceUsage: true,
	RunE:         runTUI,
}

var (
	flagHost     string
	flagPort     string
	flagPath     string
	flagPassword string
	flagCert     string
	flagInterval time.Duration
	flagTheme    string
)

func init() {
	cfg := i2ptui.LoadConfig(i2ptui.DefaultConfigPath())

	rootCmd.Flags().StringVar(&flagHost, "host", defaultStr(cfg.Host, "127.0.0.1"), "I2PControl host")
	rootCmd.Flags().StringVar(&flagPort, "port", defaultStr(cfg.Port, "7657"), "I2PControl port")
	rootCmd.Flags().StringVar(&flagPath, "path", defaultStr(cfg.Path, "jsonrpc"), "RPC URL path")
	rootCmd.Flags().StringVar(&flagPassword, "password", defaultStr(cfg.Password, "itoopie"), "I2PControl API password")
	rootCmd.Flags().StringVar(&flagCert, "cert", cfg.Cert, "path to self-signed cert")
	rootCmd.Flags().DurationVar(&flagInterval, "interval", 5*time.Second, "polling interval")
	rootCmd.Flags().StringVar(&flagTheme, "theme", defaultStr(cfg.Theme, "dark"), "color theme (dark, light)")

	rootCmd.AddCommand(versionCmd)
}

// versionCmd prints the build version of i2ptui.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of i2ptui",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("i2ptui", version)
	},
}

// runTUI launches the Bubble Tea program with the configured options.
func runTUI(cmd *cobra.Command, args []string) error {
	var opts []i2ptui.Option
	opts = append(opts,
		i2ptui.WithHost(flagHost),
		i2ptui.WithPort(flagPort),
		i2ptui.WithPath(flagPath),
		i2ptui.WithPassword(flagPassword),
		i2ptui.WithCert(flagCert),
		i2ptui.WithInterval(flagInterval),
	)
	if flagTheme == "light" {
		opts = append(opts, i2ptui.WithTheme(i2ptui.LightTheme))
	}

	m := i2ptui.New(opts...)

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("runtime error: %w", err)
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// defaultStr returns val if non-empty, otherwise the given fallback value.
func defaultStr(val, fallback string) string {
	if val != "" {
		return val
	}
	return fallback
}
