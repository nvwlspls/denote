package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {

	// check for config env var

	// load config

	// make helpful command line output (ie denote -help)

	// check for command line input to see if we should make a file with today's date
	// or one from the given arg

	// set the denote directory
	// check for custom directory for config
	// if custom directory is not set then use the default
	denoteDir := os.Getenv("HOME") + "/.denote/"
	denoteDirError := ensureDenoteDirExists(denoteDir)
	if denoteDirError != nil {
		fmt.Printf("Error ensuring denote directory exists: %v\n", denoteDirError)
		os.Exit(1)
	}
	// get current date
	currentDate := getCurrentDate()
	currentDateFileName := currentDate + ".md"

	openEditor(currentDateFileName, denoteDir)
}

func getCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func openEditor(fileName string, dir string) {

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default to vim if EDITOR is not set
	}

	cmd := exec.Command(editor, "+normal Go", dir+fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error opening editor: %v\n", err)
	}
}

func ensureDenoteDirExists(denoteDir string) error {
	expandedDir, err := expandPath(denoteDir)
	if err != nil {
		return fmt.Errorf("failed to expand path: %v", err)
	}

	if _, err := os.Stat(expandedDir); os.IsNotExist(err) {
		if err := os.MkdirAll(expandedDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		fmt.Printf("Directory created: %s\n", expandedDir)
	} else if err != nil {
		return fmt.Errorf("error checking directory: %v", err)
	} else {
		fmt.Printf("Directory already exists: %s\n", expandedDir)
	}

	return nil
}

func expandPath(path string) (string, error) {
	if path[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return homeDir + path[1:], nil
	}
	return path, nil
}
