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
	paths, err := makePaths()
	if err != nil {
		handleFatalError(err)
	}

	// Initialize paths
	err = InitPaths(paths)
	if err != nil {
		handleFatalError(err)
	}

	// Load environment variables from .env file
	err = godotenv.Load(".env")
	if err != nil {
		err := os.WriteFile(".env", []byte("OPENAI_API_KEY = \"\""), 0775)
		if err != nil {
			handleFatalError(fmt.Errorf("no \".env\" file was located, its is needed to keep your OPENAI_API key in this file"))
		}
	}

	// Check if OPENAI_API_KEY is set and not empty
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		handleFatalError(fmt.Errorf("OPENAI_API_KEY is not set or is empty, please check the \".env\" file"))
	}

	// Retrieve PDF file paths
	pdfs, err := getPDFS(paths)
	if err != nil {
		handleFatalError(err)
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
	for idx, pdf := range pdfs {
		pdfsList.AddItem(filepath.Base(pdf), "", rune(49+idx), func() {
			showLoadingSpinner(pages, pdf, paths, finalMessage)
		})
	}

	// Start the tview application
	if err := app.SetRoot(pages, true).Run(); err != nil {
		handleFatalError(err)
	}
}

// showLoadingSpinner displays the loading spinner and processes the selected PDF
func showLoadingSpinner(pages *tview.Pages, pdfPath string, paths map[string]string, finalMessage *tview.TextView) {
	// Clear the screen before switching to the loading page
	app.Sync()

	pages.SwitchToPage("loading")

	go func() {
		err := ProcessPDF(pdfPath, paths["images"])
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

		err = GenerateQuestions(paths["images"], paths["questions"], pdfPath)
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

func makePaths() (map[string]string, error) {
	var paths = make(map[string]string)

	wd, err := os.Getwd()
	if err != nil {
		return paths, err
	}

	paths["data"] = filepath.Join(wd, "data")
	paths["pdfs"] = filepath.Join(paths["data"], "pdfs")
	paths["images"] = filepath.Join(paths["data"], "images")
	paths["questions"] = filepath.Join(paths["data"], "questions")

	return paths, nil
}

func getPDFS(paths map[string]string) ([]string, error) {

	var pdfs []string
	err := filepath.Walk(paths["pdfs"], func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".pdf" {
			pdfs = append(pdfs, path)
		}
		return nil
	})
	if err != nil {
		return pdfs, err
	}

	if len(pdfs) == 0 {
		return pdfs, fmt.Errorf("no PDFs found in %q folder", paths["pdfs"])
	}

	return pdfs, nil
}

// handleFatalError handles fatal errors by printing them to the terminal and waiting for a keypress before exiting
func handleFatalError(err error) {
	fmt.Printf("Error: %v\n", err)
	fmt.Println("Press any key to exit...")
	fmt.Scanln()
	os.Exit(1)
}
