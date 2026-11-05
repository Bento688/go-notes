package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const NotesFileName = "NOTES.md"
const ContentTemplate = `# {{.Date}} 

--- 

## TODAY'S TASKS 
-

--- 

## IMPORTANT NOTES 
-

--- 

## THINGS DONE TODAY 
-`

type NotesData struct {
	Date string
}

func main() {

	// get current time
	now := time.Now()

	// e.g. 01-09-2026
	// date for the notes directory.
	notesDirDate := now.Format("02-01-2006")

	// e.g. MONDAY, SEPTEMBER 1, 2026
	// date for inside the note.
	notesHeaderDate := strings.ToUpper(now.Format("Monday, January 2, 2006"))
	fmt.Println(notesHeaderDate)

	// 2026-september
	// date for the month level directory
	monthDirDate := strings.ToLower(now.Format("2006-January"))

	absoluteDayDir := filepath.Join(
		"/Users/benedictestefanhalkin/Developer/lndata/md-notes/carbon64",
		monthDirDate,
		notesDirDate,
	)

	fmt.Printf("Target Directory: %s\n", absoluteDayDir)

	if isDirectory(absoluteDayDir) {
		// assert if NOTES.md exist
		if assertNotesExist(absoluteDayDir) {
			fmt.Println("[SKIP] NOTES.md already exists, skipping file creation.")
		} else {
			fmt.Println("[FOUND] Target directory exists, creating notes file...")
			// if not exist, then create
			err := createNotesFile(absoluteDayDir, notesHeaderDate)
			if err != nil {
				fmt.Println("[ERROR] Error during NOTES.md creation! Aborting")
				os.Exit(1)
			}
			fmt.Println("[SUCCESS] NOTES.md created with template.")
		}

	} else {
		// absolute day dir
		// does not exist, create it
		err := os.MkdirAll(absoluteDayDir, 0755)
		if err != nil {
			fmt.Println("[ERROR] Error during NOTES.md creation! Aborting")
			os.Exit(1)
		}
		fmt.Println("[SUCCESS] Directory created successfully! Creating NOTES.md")

		// create the directory, then create the notes
		err2 := createNotesFile(absoluteDayDir, notesHeaderDate)
		if err2 != nil {
			fmt.Println("[ERROR] Error during NOTES.md creation! Aborting")
			os.Exit(1)
		}
		fmt.Println("[SUCCESS] NOTES.md created with template.")
	}

}

func createNotesFile(absoluteDayDir string, notesHeaderDate string) error {
	notesPath := filepath.Join(absoluteDayDir, NotesFileName)

	tmpl, err := template.New("notes").Parse(ContentTemplate)
	if err != nil {
		return err
	}

	var content strings.Builder

	err = tmpl.Execute(&content, NotesData{
		Date: notesHeaderDate,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(notesPath, []byte(content.String()), 0644)
}

func isDirectory(absoluteDayDir string) bool {
	fileInfo, err := os.Stat(absoluteDayDir)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func assertNotesExist(absoluteDayDir string) bool {
	notesPath := filepath.Join(absoluteDayDir, NotesFileName)

	_, err := os.Stat(notesPath)
	if err == nil {
		return true // file exists
	}

	if errors.Is(err, os.ErrNotExist) {
		return false // file does not exist
	}

	// file might exist but couldn't access, return true for safety
	return true
}
