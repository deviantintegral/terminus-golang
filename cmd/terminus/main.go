// Package main is the entry point for the terminus CLI.
package main

import (
	"fmt"
	"os"

	"github.com/deviantintegral/terminus-golang/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
