package player

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Drack112/go-youtube/pkg/logger"
)

type PlayerType string

const (
	PlayerMPV PlayerType = "mpv"
)

func DetectAvailablePlayer() (PlayerType, error) {
	if _, err := exec.LookPath("mpv"); err == nil {
		logger.Debug("[Player] Found mpv")
		return PlayerMPV, nil
	}
	return "", fmt.Errorf("mpv not found - please install mpv media player")
}

func DetectYtDlp() string {
	if _, err := exec.LookPath("yt-dlp"); err == nil {
		logger.Debug("[Player] Found yt-dlp")
		return "yt-dlp"
	}
	if _, err := exec.LookPath("youtube-dl"); err == nil {
		logger.Debug("[Player] Found youtube-dl")
		return "youtube-dl"
	}
	logger.Warn("[Player] No yt-dlp or youtube-dl found - playback may fail")
	return ""
}

func StreamVideo(videoURL string, playerType PlayerType, quality string, windowMode string) error {
	logger.Debug("[Player] Streaming video", "url", videoURL, "quality", quality, "window", windowMode)
	var logFile *os.File
	if logger.LogFile != nil {
		logFile = logger.LogFile
	} else {
		lf, logPath, err := createLogFile()
		if err != nil {
			logger.Warn("[Player] Failed to create log file", "error", err)
		} else {
			defer lf.Close()
			logFile = lf

			if err := openLogWindow(logPath); err != nil {
				logger.Warn("[Player] Failed to open log window", "error", err)
			}
		}
	}

	cmd := buildMPVCommandWithOptions(videoURL, quality, windowMode)
	// keep reference to allow external shutdown
	currentMPVCmd = cmd

	if logFile != nil {
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	logger.Debug("[Player] Executing command", "cmd", cmd.String())

	if err := cmd.Run(); err != nil {
		logger.Error("[Player] Failed to run player", "error", err)
		// clear currentMPVCmd even on error
		currentMPVCmd = nil
		// close tail window if any
		logger.CloseTailWindow()
		return fmt.Errorf("failed to run mpv: %w", err)
	}

	// mpv exited normally - clear reference and close tail
	currentMPVCmd = nil
	logger.CloseTailWindow()

	return nil
}

// currentMPVCmd holds the running mpv command so it can be stopped externally
var currentMPVCmd *exec.Cmd

// StopCurrentPlayer attempts to gracefully stop the running mpv process
func StopCurrentPlayer() error {
	if currentMPVCmd == nil || currentMPVCmd.Process == nil {
		return nil
	}

	if err := currentMPVCmd.Process.Signal(syscall.SIGTERM); err != nil {
		// fallback to Kill
		_ = currentMPVCmd.Process.Kill()
	}

	currentMPVCmd = nil
	logger.CloseTailWindow()
	return nil
}

func buildMPVCommandWithOptions(videoURL string, quality string, windowMode string) *exec.Cmd {
	ytdlp := DetectYtDlp()

	args := []string{
		"--osd-level=1",
		"--osd-duration=2000",
		"--osd-status-msg=${time-pos} / ${duration}",
	}

	switch strings.ToLower(windowMode) {
	case "fullscreen", "fs":
		args = append(args, "--fullscreen")
	case "windowed", "window", "":
		args = append(args, "--force-window=yes")
	case "borderless":
		args = append(args, "--force-window=yes", "--no-border")
	case "maximized", "max":
		args = append(args, "--force-window=yes", "--window-maximized")
	default:
		args = append(args, "--force-window=yes")
	}

	if ytdlp != "" {
		args = append(args, "--script-opts=ytdl_hook-ytdl_path="+ytdlp)

		userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
		args = append(args, "--user-agent="+userAgent)

		ytdlFormat := convertQualityToFormat(quality)
		args = append(args, fmt.Sprintf("--ytdl-format=%s", ytdlFormat))
	} else {
		args = append(args, "--ytdl-format=best")
		logger.Warn("[Player] yt-dlp not found - playback will likely fail. Install with: pip install yt-dlp")
	}

	args = append(args, videoURL)
	return exec.Command("mpv", args...)
}

func convertQualityToFormat(quality string) string {
	quality = strings.ToLower(strings.TrimSpace(quality))

	switch quality {
	case "best", "highest", "max":
		return "bestvideo+bestaudio/best"
	case "worst", "lowest", "min":
		return "worstvideo+worstaudio/worst"
	case "1080p", "1080", "hd":
		return "bestvideo[height<=1080]+bestaudio/best[height<=1080]"
	case "720p", "720":
		return "bestvideo[height<=720]+bestaudio/best[height<=720]"
	case "480p", "480":
		return "bestvideo[height<=480]+bestaudio/best[height<=480]"
	case "360p", "360":
		return "bestvideo[height<=360]+bestaudio/best[height<=360]"
	case "audio", "audio-only":
		return "bestaudio"
	default:
		return "bestvideo+bestaudio/best"
	}
}

func DownloadVideo(videoURL string, container string, quality string, withThumb bool) error {
	ytdlp := DetectYtDlp()
	if ytdlp == "" {
		return fmt.Errorf("yt-dlp or youtube-dl not found; install yt-dlp to enable downloads")
	}

	args := []string{"-o", "%(title)s.%(ext)s"}

	if withThumb {
		args = append(args, "--write-thumbnail")
	}

	if quality != "" {
		qfmt := convertQualityToFormat(quality)
		args = append(args, "-f", qfmt)
	}

	if container != "" {
		args = append(args, "--recode-video", container)
	}

	args = append(args, videoURL)

	cmd := exec.Command(ytdlp, args...)

	var out io.Writer
	if logger.LogFile != nil {
		out = logger.LogFile
	} else {
		out = io.Discard
	}
	cmd.Stdout = out
	cmd.Stderr = out

	return cmd.Run()
}
