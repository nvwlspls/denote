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
	lastFile := findLastFile(denoteDir, currentDateFileName)
	fmt.Printf("the last file is %s\n", lastFile)
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
	// TODO: Get some logging going
	fmt.Printf("Your editory is %s\n", editor)
	cmd := &exec.Cmd{}
	initDailyFile(fileName, dir)
	if editor == "vim" {
		cmd = exec.Command(editor, "+normal G", dir+fileName)
	} else {
		cmd = exec.Command(editor, dir+fileName)
	}
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

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func initDailyFile(fileName string, dir string) {
	// if the file doesn't exists we will need to create it in order to add the
	// first line date
	if !fileExists(dir + fileName) {
		// Create the file and add the current date as the first line
		file, err := os.Create(dir + fileName)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
		defer file.Close()

		// Write the current date as the first line
		_, err = file.WriteString("# " + getCurrentDate() + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
}

// parse todos
// write a function that will look at the most recent file before today and
// carry over the todos

// first find the latest file
func findLastFile(denoteDir string, todaysFile string) string {
	files, err := os.ReadDir(denoteDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return ""
	}

	var latestFile string
	var latestTime time.Time

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileTime, err := time.Parse("2006-01-02.md", file.Name())
		if err != nil {
			continue // Skip files that don't match the date format
		}

		if fileTime.Before(today) && fileTime != today && file.Name() != todaysFile && (latestFile == "" || fileTime.After(latestTime)) {
			latestFile = file.Name()
			latestTime = fileTime
		}
	}

	return latestFile
}
