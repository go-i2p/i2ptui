# i2ptui

Embeddable TUI and freestanding CLI for monitoring and managing an I2P router
via I2PControl-RPC. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
and [go-i2p/go-i2pcontrol](https://github.com/go-i2p/go-i2pcontrol).

## Features

- **Dashboard** — live view of router status, bandwidth, known peers, participating tunnels, and uptime
- **Stats** — detailed bandwidth rates, tunnel build success/reject/expire percentages, and build request times
- **Peers** — peer count overview with reseeding and net-status indicators
- **Control** — restart, graceful restart, shutdown, graceful shutdown, and update checks with confirmation prompts
- **Settings** — edit bandwidth limits and share percentage via I2PControl with confirmation dialog and restart-required indicator
- **Live Graphs** — Unicode sparklines for bandwidth, tunnels, and peer count; horizontal bar chart for build success rates; toggle with `g`
- **Notifications** — in-TUI notification bar for status changes with optional desktop notifications (`notify-send` / `osascript`)
- **Theming** — built-in dark and light color themes selectable via `--theme` flag or config file
- **Mouse support** — click tabs to switch; works alongside keyboard navigation
- **Persistent config** — defaults stored in `~/.config/i2ptui/config.json`
- **Embeddable** — the root `tea.Model` is exported so other Go programs can host the TUI inside their own Bubble Tea application
- **Single binary** — `go build ./cmd/i2ptui` produces a static binary with no runtime dependencies beyond a reachable I2PControl port

## Install

```sh
go install github.com/go-i2p/i2ptui/cmd/i2ptui@latest
```

## Usage

```sh
i2ptui [flags]
i2ptui [command]
```

### Commands

| Command | Description |
|---------|-------------|
| `completion` | Generate shell completion scripts (bash, zsh, fish, powershell) |
| `version` | Print the build version |
| `help` | Help about any command |

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `127.0.0.1` | I2PControl host |
| `--port` | `7657` | I2PControl port |
| `--path` | `jsonrpc` | RPC URL path |
| `--password` | `itoopie` | I2PControl API password |
| `--cert` | | Path to self-signed cert (enables TLS skip-verify) |
| `--interval` | `5s` | Polling interval |
| `--theme` | `dark` | Color theme (`dark`, `light`) |

### Shell Completions

```sh
# Bash
i2ptui completion bash > /etc/bash_completion.d/i2ptui

# Zsh
i2ptui completion zsh > "${fpath[1]}/_i2ptui"

# Fish
i2ptui completion fish > ~/.config/fish/completions/i2ptui.fish
```

## Embedding

```go
package main

import (
    "github.com/go-i2p/i2ptui"
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    m := i2ptui.New(
        i2ptui.WithHost("127.0.0.1"),
        i2ptui.WithPort("7657"),
        i2ptui.WithPassword("itoopie"),
    )
    p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
    p.Run()
}
```

## Key Bindings

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Cycle tabs |
| `1` `2` `3` `4` `5` | Jump to Dashboard, Stats, Peers, Control, Settings |
| `↑`/`k` `↓`/`j` | Navigate |
| `Enter` | Activate |
| `Esc` | Cancel / dismiss |
| `g` | Toggle graph panel |
| `r` | Force refresh |
| `q` / `Ctrl+C` | Quit |

## License

MIT
