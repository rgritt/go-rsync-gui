package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	sourcePath     string
	destPath       string
	cmd            *exec.Cmd
	transferActive bool
	mu             sync.Mutex
	logFile        *os.File
)

func main() {
	// Set up logging to a file
	var err error
	logFile, err = os.OpenFile("rsync_transfer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("Application started")

	// Create a new Fyne app
	a := app.New()
	w := a.NewWindow("Rsync GUI")

	// Create labels to show selected paths
	sourceLabel := widget.NewLabel("Source Path: Not selected")
	destLabel := widget.NewLabel("Dest Path: Not selected")

	// Button to select source path
	sourceButton := widget.NewButton("Select Source Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				sourcePath = uri.Path()
				sourceLabel.SetText(fmt.Sprintf("Source Path: %s", sourcePath))
				log.Println("Selected source path:", sourcePath)
			}
		}, w)
	})

	// Button to select destination path
	destButton := widget.NewButton("Select Destination Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				destPath = uri.Path()
				destLabel.SetText(fmt.Sprintf("Dest Path: %s", destPath))
				log.Println("Selected destination path:", destPath)
			}
		}, w)
	})

	// Button to start rsync process
	rsyncButton := widget.NewButton("Start Rsync", func() {
		if sourcePath == "" || destPath == "" {
			dialog.ShowInformation("Error", "Please select both source and destination paths", w)
			log.Println("Error: Source or destination path not selected")
			return
		}

		if transferActive {
			dialog.ShowInformation("Error", "A transfer is already in progress", w)
			log.Println("Error: Transfer already in progress")
			return
		}

		// Run rsync with sourcePath and destPath in a goroutine
		go func() {
			mu.Lock()
			transferActive = true
			mu.Unlock()

			log.Println("Starting rsync transfer")
			err := runRsync(sourcePath, destPath)
			mu.Lock()
			transferActive = false
			mu.Unlock()

			if err != nil {
				dialog.ShowInformation("Error", fmt.Sprintf("Rsync failed: %s", err.Error()), w)
				log.Println("Error: Rsync failed:", err)
			} else {
				dialog.ShowInformation("Success", "Rsync completed successfully", w)
				log.Println("Rsync completed successfully")
			}
		}()
	})

	// Button to stop rsync process
	stopButton := widget.NewButton("Stop Rsync", func() {
		if cmd != nil && transferActive {
			err := cmd.Process.Kill()
			if err != nil {
				dialog.ShowInformation("Error", fmt.Sprintf("Failed to stop rsync: %s", err.Error()), w)
				log.Println("Error: Failed to stop rsync:", err)
			} else {
				dialog.ShowInformation("Stopped", "Rsync process has been stopped", w)
				log.Println("Rsync process stopped")
			}
		} else {
			dialog.ShowInformation("Error", "No active rsync process to stop", w)
			log.Println("Error: No active rsync process to stop")
		}
	})

	// Create a vertical box layout with all the buttons and labels
	content := container.NewVBox(
		sourceLabel,
		sourceButton,
		destLabel,
		destButton,
		rsyncButton,
		stopButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 200))
	w.ShowAndRun()

	log.Println("Application closed")
}

// Function to run the rsync command
func runRsync(source, dest string) error {
	// Construct the rsync command
	cmd = exec.Command("rsync", "-avh", source+string(filepath.Separator), dest)

	// Set the output to the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Log the command execution
	log.Printf("Executing command: rsync -avh %s %s\n", source, dest)

	// Run the command and return any errors
	return cmd.Run()
}
