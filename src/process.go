package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProcessPDF processes each PDF file by converting it to images
func ProcessPDF(pdfPath, imagesDir string) error {
	// Create a folder for images of the current PDF
	baseName := strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))
	folder := filepath.Join(imagesDir, baseName)

	// Check if folder exists, if not create it
	if _, err := os.Stat(folder); err != nil {
		if os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			err = os.Mkdir(folder, 0775)
			if err != nil {
				return fmt.Errorf("failed to create directory %q: %v", folder, err)
			}
			fmt.Printf("Created directory %q\n", folder)
		} else {
			// Return error if it is not a "not exist" error
			return fmt.Errorf("failed to check directory %q: %v", folder, err)
		}
	} else {
		fmt.Printf("Folder %q already exists.\n", folder)
	}

	// Convert PDF to images using pdftocairo
	outputPattern := filepath.Join(folder, baseName)
	cmd := exec.Command("pdftocairo", "-jpeg", pdfPath, outputPattern)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to convert PDF to images: %v, output: %s", err, output)
	}

	return nil
}