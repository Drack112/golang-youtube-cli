package logger

import (
	"fmt"
	stdlog "log"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/lipgloss"
	clog "github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

var Logger *clog.Logger
var LogFile *os.File
var LogFilePath string

func getColoredPrefix() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#6366F1")).
		Bold(true).Padding(0, 1).
		MarginRight(1)

	return style.Render("Go-YouTube")
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

	Logger.SetColorProfile(termenv.EnvColorProfile())

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
