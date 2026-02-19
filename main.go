package main

import (
	"fmt"
	"os"

	"github.com/tristian-me/chksum/checker"
)

func main() {
	args := os.Args[1:]

	// Check directory
	if len(args) == 0 {
		if err := checker.CheckDir(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return
	}

	for _, file := range args {
		if err := checker.CheckFile(file); err != nil {
			fmt.Println(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}
	}
}
