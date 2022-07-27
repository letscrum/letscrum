package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/letscrum/letscrum/tools/releasenoter/note"
	"github.com/letscrum/letscrum/tools/releasenoter/utils"
)

func main() {
	var notesDir, outPath string

	flag.StringVar(&notesDir, "notesDir", "./changelogs/0.1.1", "contain yaml directory")
	flag.StringVar(&outPath, "outPath", "./changelogs/tool/release_note", "out file path")
	flag.Parse()

	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		fmt.Printf("not find notes directory, %s not exist.\n", notesDir)
		os.Exit(1)
	}
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		fmt.Printf("not find out path, %s not exist.\n", outPath)
		os.Exit(1)
	}

	version := utils.GetVersionName(notesDir)
	notesFiles := utils.GetDirFiles(notesDir)
	if len(notesFiles) < 1 {
		fmt.Println("notes num is zero")
		os.Exit(1)
	}

	notes, err := note.ParseNotesFile(notesFiles)
	if err != nil {
		fmt.Println("parse notes err:", err.Error())
		os.Exit(1)
	}

	if err := note.CreateMarkDown(notes, outPath, version); err != nil {
		fmt.Println("create markdown file err:", err.Error())
		return
	}

	fmt.Printf("note file %s out success\n", version)
}
