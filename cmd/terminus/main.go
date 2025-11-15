// Package main is the entry point for the terminus CLI.
package main

import (
	"os"

	"github.com/pantheon-systems/terminus-go/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
