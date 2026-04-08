package i2ptui

import (
	"runtime"
	"testing"
)

func TestNotifyCmdLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("linux-only test")
	}
	cmd := notifyCmd("title", "body")
	if cmd == nil {
		t.Fatal("expected non-nil cmd on linux")
	}
	args := cmd.Args
	if len(args) != 3 || args[0] != "notify-send" || args[1] != "title" || args[2] != "body" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestNotifyCmdDarwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-only test")
	}
	cmd := notifyCmd("title", "body")
	if cmd == nil {
		t.Fatal("expected non-nil cmd on darwin")
	}
	if cmd.Args[0] != "osascript" {
		t.Fatalf("expected osascript, got %s", cmd.Args[0])
	}
}

func TestNotifyCmdQuoteEscaping(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-only test")
	}
	cmd := notifyCmd(`he"llo`, `wor"ld`)
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}
	script := cmd.Args[2]
	if testing.Verbose() {
		t.Logf("script: %s", script)
	}
	// Verify the double quotes in input were escaped.
	if script == `display notification "wor"ld" with title "he"llo"` {
		t.Fatal("quotes were not escaped")
	}
}

func TestNotifyCmdCurrentOS(t *testing.T) {
	cmd := notifyCmd("test", "msg")
	switch runtime.GOOS {
	case "linux", "darwin":
		if cmd == nil {
			t.Fatal("expected non-nil cmd on supported OS")
		}
	default:
		if cmd != nil {
			t.Fatalf("expected nil cmd on unsupported OS %s", runtime.GOOS)
		}
	}
}
