package i2ptui

import (
	"os/exec"
	"runtime"
	"strings"
)

// notifyCmd builds the *exec.Cmd for an OS desktop notification.
// Returns nil when the current OS is unsupported.
func notifyCmd(title, body string) *exec.Cmd {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("notify-send", title, body)
	case "darwin":
		safeBody := strings.ReplaceAll(body, `"`, `\"`)
		safeTitle := strings.ReplaceAll(title, `"`, `\"`)
		script := `display notification "` + safeBody + `" with title "` + safeTitle + `"`
		return exec.Command("osascript", "-e", script)
	default:
		return nil
	}
}

// desktopNotify sends an OS desktop notification if supported.
func desktopNotify(title, body string) {
	if cmd := notifyCmd(title, body); cmd != nil {
		_ = cmd.Start()
	}
}
