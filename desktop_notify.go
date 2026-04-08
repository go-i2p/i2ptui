package i2ptui

import (
	"os/exec"
	"runtime"
	"strings"
)

// desktopNotify sends an OS desktop notification if supported.
func desktopNotify(title, body string) {
	switch runtime.GOOS {
	case "linux":
		// Best-effort; ignore errors if notify-send is unavailable.
		_ = exec.Command("notify-send", title, body).Start()
	case "darwin":
		// Escape quotes to prevent injection in osascript.
		safeBody := strings.ReplaceAll(body, `"`, `\"`)
		safeTitle := strings.ReplaceAll(title, `"`, `\"`)
		script := `display notification "` + safeBody + `" with title "` + safeTitle + `"`
		_ = exec.Command("osascript", "-e", script).Start()
	}
}
