package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

var app *tview.Application = tview.NewApplication()

func main() {
	// Initialize paths
	_, _, pdfsDir, imagesDir, questionsDir, err := InitPaths()
	if err != nil {
		handleFatalError(err)
	}

	// Load environment variables from .env file
	err = godotenv.Load(".env")
	if err != nil {
		handleFatalError(fmt.Errorf("no \".env\" file was located, its is needed to keep your OPENAI_API key in this file"))
	}

	// Check if OPENAI_API_KEY is set and not empty
	apiKey := os.Getenv("OPENAI_API")
	if apiKey == "" {
		handleFatalError(fmt.Errorf("OPENAI_API is not set or is empty, please check the \".env\" file"))
	}

	// Retrieve PDF file paths
	var pdfPaths []string
	err = filepath.Walk(pdfsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".pdf" {
			pdfPaths = append(pdfPaths, path)
		}
		return nil
	})
	if err != nil {
		handleFatalError(err)
	}

	if len(pdfPaths) == 0 {
		handleFatalError(fmt.Errorf("no PDFs found in %q folder", pdfsDir))
	}

	// Set up tview components
	pdfsList := tview.NewList()
	listFlx := Center(0, 0, "Select PDF File!", pdfsList)

	loading := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	go animateLoadingSpinner(loading)

	loadingFlx := Center(0, 0, "Generating Questions!", loading)

	finalMessage := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	finalMessageFlx := Center(0, 0, "Operations Complete!", finalMessage)

	pages := tview.NewPages().
		AddPage("list", listFlx, true, true).
		AddPage("loading", loadingFlx, true, false).
		AddPage("final", finalMessageFlx, true, false)

	// Add PDFs to the list view
	for idx, pdfPath := range pdfPaths {
		pdfPath := pdfPath // capture range variable
		pdfsList.AddItem(filepath.Base(pdfPath), "", rune(49+idx), func() {
			showLoadingSpinner(pages, pdfPath, imagesDir, questionsDir, finalMessage)
		})
	}

	// Start the tview application
	if err := app.SetRoot(pages, true).Run(); err != nil {
		handleFatalError(err)
	}
}

// handleFatalError handles fatal errors by printing them to the terminal and waiting for a keypress before exiting
func handleFatalError(err error) {
	fmt.Printf("Error: %v\n", err)
	fmt.Println("Press any key to exit...")
	fmt.Scanln()
	os.Exit(1)
}

// showLoadingSpinner displays the loading spinner and processes the selected PDF
func showLoadingSpinner(pages *tview.Pages, pdfPath string, imagesDir string, questionsDir string, finalMessage *tview.TextView) {
	// Clear the screen before switching to the loading page
	app.Sync()

	pages.SwitchToPage("loading")

	go func() {
		err := ProcessPDF(pdfPath, imagesDir)
		if err != nil {
			app.QueueUpdateDraw(func() {
				finalMessage.SetText(fmt.Sprintf("[red]%v", err))
				pages.SwitchToPage("final")
			})
			// Delay to show the error message before navigating back to the list
			time.Sleep(3 * time.Second)
			app.QueueUpdateDraw(func() {
				pages.SwitchToPage("list")
			})
			return
		}

		err = GenerateQuestions(imagesDir, questionsDir, pdfPath)
		if err != nil {
			app.QueueUpdateDraw(func() {
				finalMessage.SetText(fmt.Sprintf("[red]%v", err))
				pages.SwitchToPage("final")
			})
			// Delay to show the error message before navigating back to the list
			time.Sleep(3 * time.Second)
			app.QueueUpdateDraw(func() {
				pages.SwitchToPage("list")
			})
			return
		}

		app.QueueUpdateDraw(func() {
			finalMessage.SetText("[green]Process Finished!")
			pages.SwitchToPage("final")
		})

		// Delay to show the success message before navigating back to the list
		time.Sleep(2 * time.Second)

		app.QueueUpdateDraw(func() {
			pages.SwitchToPage("list")
		})
	}()
}

// animateLoadingSpinner animates the loading spinner in the loading text view
func animateLoadingSpinner(loading *tview.TextView) {
	spinner := []string{"|", "/", "-", "\\"}
	for {
		for _, frame := range spinner {
			app.QueueUpdateDraw(func() {
				loading.SetText(frame + " Generating...")
			})
			time.Sleep(100 * time.Millisecond)
		}
	}
}