package main

import (
	"embed"
	"os"

	"github.com/mukezhz/geng/cmd"
)

//go:embed templates/wesionary/*
var templatesFS embed.FS

func main() {
	if err := cmd.Execute(templatesFS); err != nil {
		os.Exit(1)
	}
}
