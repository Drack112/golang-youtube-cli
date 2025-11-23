package utils

import (
	"fmt"
	"os"
)

var Version = "v0.1.0"

func HasVersionArg() bool {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		return arg == "--version" || arg == "-v"
	}

	return false
}

func ShowVersion() {
	fmt.Printf("GO-Youtube v%s", Version)
}
