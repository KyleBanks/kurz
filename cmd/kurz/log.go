package main

import (
	"fmt"
	"os"
	"strings"
)

func log(msg string, a ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf(msg, a...)
}

func logError(err error) {
	log("ERROR: %v\n", err)
	os.Exit(1)
}

func printUsage(code int) {
	name := os.Args[0]

	log(`%v allows you to view markdown documents on the command-line in a feature-rich UI. 

Usage:
  %v [options] path 
    	Where 'path' is a local file, remote URL, or Git repository.

Example:
  %v ./path/to/file.md
  %v http://example.com/document.md
  %v github.com/KyleBanks/modoc

To print this message, use the '--help' flag.`, name, name, name, name, name)
	os.Exit(code)
}
