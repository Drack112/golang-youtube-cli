package main

import (
	"fmt"

	"github.com/Drack112/go-youtube/internal/flags"
	"github.com/Drack112/go-youtube/internal/handlers"
	"github.com/Drack112/go-youtube/pkg/logger"
)

func main() {
	value, err := flags.ParseFlags()
	if err != nil {
		if err == flags.ErrHelpRequested {
			return
		}
		logger.Fatal(flags.ErrorHandler(err))
	}

	fmt.Println(handlers.SearchWithRetries(value))
}
