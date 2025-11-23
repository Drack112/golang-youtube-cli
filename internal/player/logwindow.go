package player

import (
	"os"
	"path/filepath"

	"github.com/Drack112/go-youtube/pkg/logger"
)

func createLogFile() (*os.File, string, error) {
	tmpDir := os.TempDir()
	logPath := filepath.Join(tmpDir, "go-youtube.log")

	logFile, err := os.Create(logPath)
	if err != nil {
		return nil, "", err
	}

	return logFile, logPath, nil
}

func openLogWindow(logPath string) error {
	logger.TryOpenLogWindow(logPath)
	return nil
}
