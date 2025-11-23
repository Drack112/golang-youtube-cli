package flags

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/Drack112/go-youtube/pkg/utils"
)

var (
	IsDebug              bool
	ErrHelpRequested     = errors.New("help requested")
	ErrDownloadRequested = errors.New("download requested")
	ErrNoInput           = errors.New("no input provided")
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

func ParseFlags() (*Options, error) {
	opts := &Options{}

	flag.BoolVar(&opts.Debug, "debug", false, "enable debug mode")
	flag.BoolVar(&opts.Download, "download", false, "download instead of streaming")
	flag.StringVar(&opts.Quality, "quality", "best", "video quality (best, worst, 720p, 1080p)")

	flag.Usage = func() {
		fmt.Println("\ngo-youtube [OPTIONS] <url | search term | video name>")
		fmt.Println("Examples:")
		fmt.Println("  go-youtube \"funny cats\"")
		fmt.Println("  go-youtube https://youtu.be/xyz")
		fmt.Println("  go-youtube -download \"metallica live\"")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println()
	}

	help := flag.Bool("help", false, "show help message")
	version := flag.Bool("version", false, "show version information")

	flag.Parse()

	IsDebug = opts.Debug

	if *version || utils.HasVersionArg() {
		utils.ShowVersion()
		return nil, ErrHelpRequested
	}

	if *help {
		flag.Usage()
		return nil, ErrHelpRequested
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return nil, ErrNoInput
	}

	input := strings.Join(args, "")
	opts.Input = strings.TrimSpace(input)

	if utils.IsYoutubeURL(opts.Input) {
		opts.InputKind = InputYoutubeURL
	} else {
		opts.InputKind = InputSearchQuery
	}

	if opts.Download {
		return nil, ErrDownloadRequested
	}

	return opts, nil

}
