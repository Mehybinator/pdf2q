package main

import (
	"fmt"
	"os"
)

// InitPaths initializes the working directory, PDF directory, image directory, and retrieves PDF paths
func InitPaths(paths map[string]string) error {
	// Create the directories if they do not exist
	for _, v := range paths {
		if err := createDirIfNotExist(v); err != nil {
			return err
		}
	}
	return nil
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
