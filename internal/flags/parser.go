package flags

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/Drack112/go-youtube/pkg/logger"
	"github.com/Drack112/go-youtube/pkg/utils"
	"github.com/charmbracelet/huh"
)

var (
	IsDebug          bool
	ErrHelpRequested = errors.New("help requested")
	ErrNoInput       = errors.New("no input provided")
)

type InputSrc int

const (
	InputNone InputSrc = iota
	InputYoutubeURL
	InputSearchQuery
)

type Options struct {
	Debug      bool
	Quality    string
	WindowMode string

	Input     string
	InputKind InputSrc
}

func ErrorHandler(err error) string {
	if IsDebug {
		return fmt.Sprintf("%+v", err)
	}
	return fmt.Sprintf("%v (run with -debug for details)", err)
}

func getUserInput() (string, error) {
	var videoName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter video or topic bellow").Description("Type a Youtube URL or a search term").Value(&videoName).Validate(func(v string) error {
				if len(strings.TrimSpace(v)) < 3 {
					return fmt.Errorf("input must have at least 3 characters")
				}
				return nil
			}),
		),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return videoName, nil
}

func ParseFlags() (*Options, error) {
	opts := &Options{}

	debug := flag.Bool("debug", false, "enable debug mode")
	help := flag.Bool("help", false, "show help message")
	versionFlag := flag.Bool("version", false, "show version information")
	quality := flag.String("quality", "best", "video quality (best, 1080p, 720p, 480p, 360p, audio)")
	windowMode := flag.String("window", "windowed", "window mode (windowed, fullscreen, borderless, maximized)")

	flag.Usage = func() {
		fmt.Println("\ngo-youtube [OPTIONS] <url | search term>")
		fmt.Println("Examples:")
		fmt.Println("  go-youtube \"lofi chill\"")
		fmt.Println("  go-youtube -quality 720p -window fullscreen https://youtu.be/xxx")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()

	opts.Debug = *debug
	opts.Quality = *quality
	opts.WindowMode = *windowMode
	IsDebug = opts.Debug

	if opts.Debug {
		logger.InitLogger(opts.Debug)
		logger.Debug("Debug mode enabled")
	}

	if *versionFlag || utils.HasVersionArg() {
		utils.ShowVersion()
		return nil, ErrHelpRequested
	}

	if *help {
		flag.Usage()
		return nil, ErrHelpRequested
	}

	args := flag.Args()
	var input string

	if len(args) > 0 {
		input = strings.Join(args, " ")
	} else {
		in, err := getUserInput()
		if err != nil {
			return nil, err
		}
		input = utils.CleanYoutubeLink(in)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return nil, ErrNoInput
	}

	opts.Input = input

	if utils.IsYouTubeURL(input) {
		opts.InputKind = InputYoutubeURL
	} else {
		opts.InputKind = InputSearchQuery
	}

	return opts, nil
}
