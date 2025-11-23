package main

import (
	"fmt"
	"log"

	"github.com/Drack112/go-youtube/internal/flags"
)

func main() {
	options, err := flags.ParseFlags()
	if err != nil {
		if err == flags.ErrDownloadRequested {
			// Handle After
			return
		}

		if err == flags.ErrHelpRequested {
			return
		}

		log.Fatalln(flags.ErrorHandler(err))
	}

	fmt.Println(options)

}
