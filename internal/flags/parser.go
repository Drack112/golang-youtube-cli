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
	IsDebug              bool
	ErrHelpRequested     = errors.New("help requested")
	ErrDownloadRequested = errors.New("download requested")
	ErrNoInput           = errors.New("no input provided")

	GlobalQuality string
)

type InputSrc int

const (
	InputNone InputSrc = iota
	InputYoutubeURL
	InputSearchQuery
)

type Options struct {
	Debug    bool
	Download bool
	Quality  string

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
	download := flag.Bool("download", false, "download instead of stream")
	quality := flag.String("quality", "best", "video quality (best, worst, 720p, 1080p)")

	help := flag.Bool("help", false, "show help message")
	versionFlag := flag.Bool("version", false, "show version information")

	flag.Usage = func() {
		fmt.Println("\ngo-youtube [OPTIONS] <url | search term>")
		fmt.Println("Examples:")
		fmt.Println("  go-youtube \"lofi chill\"")
		fmt.Println("  go-youtube https://youtu.be/xxx")
		fmt.Println("  go-youtube -download \"metallica live\"")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()

	opts.Debug = *debug
	opts.Download = *download
	opts.Quality = *quality

	IsDebug = opts.Debug
	GlobalQuality = opts.Quality

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

	if opts.Download {
		return opts, ErrDownloadRequested
	}

	return opts, nil
}
