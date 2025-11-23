package player

import (
	"fmt"
	"os"
	"os/exec"
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
	var cmd *exec.Cmd

	if _, err := exec.LookPath("mate-terminal"); err == nil {
		cmd = exec.Command("mate-terminal", "--", "tail", "-f", logPath)
	} else if _, err := exec.LookPath("xterm"); err == nil {
		cmd = exec.Command("xterm", "-e", fmt.Sprintf("tail -f %s", logPath))
	} else {
		return fmt.Errorf("no supported terminal emulator found")
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// store reference so other packages can close it
	logger.TailCmd = cmd

	return nil
}
