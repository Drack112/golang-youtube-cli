package main

import (
	"fmt"
	"os"

	"github.com/Drack112/go-youtube/internal/flags"
	"github.com/Drack112/go-youtube/internal/tui"
	"github.com/Drack112/go-youtube/pkg/logger"
)

func main() {
	opts, err := flags.ParseFlags()
	if err != nil {
		switch err {
		case flags.ErrHelpRequested:
			return
		default:
			logger.Fatal(flags.ErrorHandler(err))
		}
	}

	model := tui.NewModel(opts)
	if err := tui.NewProgram(model).Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
