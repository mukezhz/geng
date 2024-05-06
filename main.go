package main

import (
	"os"

	"github.com/mukezhz/geng/cmd"
)

func main() {
  // execute the root command
	if err := cmd.Root.Execute(); err != nil {
		os.Exit(1)
	}
}
