package main

import (
	"github.com/pauwol/touch/cmd"
)

func main() {

	cmd := cmd.CMD{}

	cmd.AddFlag("--recursive", "-r", "Recursively create directories", "")
	cmd.AddFlag("--force", "-f", "Overwrite existing files", "")
	cmd.AddFlag("--directory", "-d", "Create directories instead of files", "")
	cmd.AddFlag("--timestamp", "-t", "Set the creation timestamp for the file (YYYY-MM-DD HH:MM:SS)", "")
	cmd.AddFlag("--permissions", "-p", "Set file permissions (USER or ADMIN)", "")
	cmd.AddFlag("--update", "-u", "Update timestamps and permissions without creating new files", "")
	cmd.AddFlag("--help", "-h", "Show help", "")

	cmd.Execute()
}
