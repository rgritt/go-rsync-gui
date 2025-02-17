package main

import (
	"fmt"
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

var sourcePath string
var destPath string
var cmd *exec.Cmd
var transferActive bool
var mu sync.Mutex

func main() {
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
			}
		}, w)
	})

	// Button to select destination path
	destButton := widget.NewButton("Select Destination Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				destPath = uri.Path()
				destLabel.SetText(fmt.Sprintf("Dest Path: %s", destPath))
			}
		}, w)
	})

	// Button to start rsync process
	rsyncButton := widget.NewButton("Start Rsync", func() {
		if sourcePath == "" || destPath == "" {
			dialog.ShowInformation("Error", "Please select both source and destination paths", w)
			return
		}

		if transferActive {
			dialog.ShowInformation("Error", "A transfer is already in progress", w)
			return
		}

		// Run rsync with sourcePath and destPath in a goroutine
		go func() {
			mu.Lock()
			transferActive = true
			mu.Unlock()

			err := runRsync(sourcePath, destPath)
			mu.Lock()
			transferActive = false
			mu.Unlock()

			if err != nil {
				dialog.ShowInformation("Error", fmt.Sprintf("Rsync failed: %s", err.Error()), w)
			} else {
				dialog.ShowInformation("Success", "Rsync completed successfully", w)
			}
		}()
	})

	// Button to stop rsync process
	stopButton := widget.NewButton("Stop Rsync", func() {
		if cmd != nil && transferActive {
			err := cmd.Process.Kill()
			if err != nil {
				dialog.ShowInformation("Error", fmt.Sprintf("Failed to stop rsync: %s", err.Error()), w)
			} else {
				dialog.ShowInformation("Stopped", "Rsync process has been stopped", w)
			}
		} else {
			dialog.ShowInformation("Error", "No active rsync process to stop", w)
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
}

// Function to run the rsync command
func runRsync(source, dest string) error {
	// Construct the rsync command
	cmd = exec.Command("rsync", "-avh", source+string(filepath.Separator), dest)

	// Set the output to the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and return any errors
	return cmd.Run()
}

