package logger

import (
	"fmt"
	stdlog "log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/lipgloss"
	clog "github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

var Logger *clog.Logger
var LogFile *os.File
var LogFilePath string
var TailCmd *exec.Cmd

func getColoredPrefix() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#6366F1")).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	return style.Render("GO-Youtube")
}

// TryOpenLogWindow attempts to open a terminal window tailing the given log path.
// It chooses a strategy depending on the OS and available terminal emulators.
func TryOpenLogWindow(path string) {
	go func() {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			if _, err := exec.LookPath("osascript"); err == nil {
				// try iTerm
				script := fmt.Sprintf(`tell application "iTerm"
	create window with default profile
	tell current session of current window to write text "tail -f %s"
end tell`, path)
				cmd = exec.Command("osascript", "-e", script)
				if err := cmd.Start(); err == nil {
					TailCmd = cmd
					return
				}

				// try Terminal.app
				script = fmt.Sprintf(`tell application "Terminal"
	do script "tail -f %s"
	activate
end tell`, path)
				cmd = exec.Command("osascript", "-e", script)
				if err := cmd.Start(); err == nil {
					TailCmd = cmd
					return
				}
			}
		case "windows":
			// Try PowerShell via wt (Windows Terminal) or powershell.exe
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

func InitLogger(isDebug bool) {
	writer := os.Stderr

	if isDebug {
		tmpDir := os.TempDir()
		logPath := filepath.Join(tmpDir, "go-youtube.log")
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err == nil {
			LogFile = f
			LogFilePath = logPath
			writer = f
			TryOpenLogWindow(logPath)
		}
	}

	Logger = clog.NewWithOptions(writer, clog.Options{
		ReportCaller:    isDebug,
		ReportTimestamp: isDebug,
		TimeFormat:      time.DateTime,
		Prefix:          getColoredPrefix(),
	})
	Logger.SetColorProfile(termenv.TrueColor)

	if isDebug {
		Logger.SetLevel(clog.DebugLevel)
		Logger.Debug("Debug logging enabled")
	} else {
		Logger.SetLevel(clog.InfoLevel)
	}

	if LogFile != nil {
		stdlog.SetOutput(LogFile)
	} else {
		stdlog.SetOutput(os.Stderr)
	}
}

func Debug(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Debug(fmt.Sprintf("%v", msg), keyvals...)
	}
}

func Info(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Info(fmt.Sprintf("%v", msg), keyvals...)
	}
}

func Warn(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Warn(fmt.Sprintf("%v", msg), keyvals...)
	}
}

func Error(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Error(fmt.Sprintf("%v", msg), keyvals...)
	}
}

func Fatal(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Fatal(fmt.Sprintf("%v", msg), keyvals...)
	}
}

func Debugf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Debug(fmt.Sprintf(format, args...))
	}
}

func Infof(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Info(fmt.Sprintf(format, args...))
	}
}

func Warnf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Warn(fmt.Sprintf(format, args...))
	}
}

func Errorf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Error(fmt.Sprintf(format, args...))
	}
}
