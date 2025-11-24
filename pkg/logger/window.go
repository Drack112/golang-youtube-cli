package logger

import (
	"fmt"
	"os/exec"
	"runtime"
)

var TailCmd *exec.Cmd

// Attempt to open a terminal window tailing the log content
func TryOpenLogWindow(path string) {
	go func() {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			// Open Terminal App on MacOS
			script := fmt.Sprintf(`tell application "Terminal"
				do script "tail -f %s"
				activate
			end tell`, path)
			cmd = exec.Command("osascript", "-e", script)
			if err := cmd.Start(); err == nil {
				TailCmd = cmd
				return
			}
		case "windows":
			if _, err := exec.LookPath("wt"); err == nil {
				cmd = exec.Command("wt", "powershell", "-NoExit", "-Command", fmt.Sprintf("Get-Content -Path '%s' -Wait", path))
			} else if _, err := exec.LookPath("powershell"); err == nil {
				cmd = exec.Command("powershell", "-NoExit", "-Command", fmt.Sprintf("Get-Content -Path '%s' -Wait", path))
			}
		default:
			// Linux/other - try common terminals
			if _, err := exec.LookPath("gnome-terminal"); err == nil {
				cmd = exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf("tail -f %s; exec bash", path))
			} else if _, err := exec.LookPath("konsole"); err == nil {
				cmd = exec.Command("konsole", "-e", "tail", "-f", path)
			} else if _, err := exec.LookPath("mate-terminal"); err == nil {
				cmd = exec.Command("mate-terminal", "--", "tail", "-f", path)
			} else if _, err := exec.LookPath("xterm"); err == nil {
				cmd = exec.Command("xterm", "-e", fmt.Sprintf("tail -f %s", path))
			}
		}

		if cmd == nil {
			return
		}

		if err := cmd.Start(); err == nil {
			TailCmd = cmd
		}
	}()
}

func CloseTailWindow() {
	if TailCmd == nil || TailCmd.Process == nil {
		return
	}

	_ = TailCmd.Process.Kill()
	TailCmd = nil
}
