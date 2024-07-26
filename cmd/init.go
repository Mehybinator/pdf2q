package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitPaths initializes the working directory, PDF directory, image directory, and retrieves PDF paths
func InitPaths() (wd string, dataDir string, pdfsDir string, imagesDir string, questionsDir string, err error) {
	// Get the current working directory
	wd, err = os.Getwd()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to get working directory: %v", err)
	}

	// Define directory paths
	dataDir = filepath.Join(wd, "data")
	pdfsDir = filepath.Join(dataDir, "pdfs")
	imagesDir = filepath.Join(dataDir, "images")
	questionsDir = filepath.Join(dataDir, "questions")

	// Create the directories if they do not exist
	dirs := []string{dataDir, pdfsDir, imagesDir, questionsDir}
	for _, dir := range dirs {
		if err := createDirIfNotExist(dir); err != nil {
			return "", "", "", "", "", err
		}
	}

	return wd, dataDir, pdfsDir, imagesDir, questionsDir, nil
}

// createDirIfNotExist creates a directory if it does not already exist
func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0775); err != nil {
			return fmt.Errorf("failed to create directory %q: %v", dir, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check directory %q: %v", dir, err)
	}
	return nil
}