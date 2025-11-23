package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

var Logger *log.Logger

func getColoredPrefix() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#6366F1")).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	return style.Render(" GO-Youtube")
}

func InitLogger(isDebug bool) {
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    isDebug,
		ReportTimestamp: isDebug,
		TimeFormat:      time.DateTime,
		Prefix:          getColoredPrefix(),
	})
	Logger.SetColorProfile(termenv.TrueColor)

	if isDebug {
		Logger.SetLevel(log.DebugLevel)
		Logger.Debug("Debug logging enabled")
	} else {
		Logger.SetLevel(log.InfoLevel)
	}
}

// Debug logs a debug message (only when debug mode is enabled)
func Debug(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Debug(fmt.Sprintf("%v", msg), keyvals...)
	}
}

// Info logs an info message
func Info(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Info(fmt.Sprintf("%v", msg), keyvals...)
	}
}

// Warn logs a warning message
func Warn(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Warn(fmt.Sprintf("%v", msg), keyvals...)
	}
}

// Error logs an error message
func Error(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Error(fmt.Sprintf("%v", msg), keyvals...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg interface{}, keyvals ...interface{}) {
	if Logger != nil {
		Logger.Fatal(fmt.Sprintf("%v", msg), keyvals...)
	}
}

// Debugf logs a formatted debug message (only when debug mode is enabled)
func Debugf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Debug(fmt.Sprintf(format, args...))
	}
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Info(fmt.Sprintf(format, args...))
	}
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Warn(fmt.Sprintf(format, args...))
	}
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Error(fmt.Sprintf(format, args...))
	}
}
