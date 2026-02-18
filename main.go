package main

import (
	"os"

	"github.com/tristian-me/chksum/checker"
)

func main() {
	args := os.Args[1:]

	// Check directory
	if len(args) == 0 {
		checker.CheckDir()
	} else {
		for _, file := range args {
			checker.CheckFile(file)
		}
	}
}
