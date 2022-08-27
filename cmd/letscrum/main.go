package main

import (
	"fmt"
	"github.com/letscrum/letscrum/internal/cmd"
	"os"
)

func main() {
	if err := cmd.GetRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
