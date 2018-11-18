package main

import (
	"github.com/rassakhatsky/tdameritrade/cmd"
)

var (
	// The app version
	// This variable is set via ld flags
	Version string
)

func main() {
	cmd.Execute()
}
