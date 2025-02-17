
# Rsync GUI

Rsync GUI is a simple graphical user interface (GUI) application built with Go and Fyne. It allows users to select source and destination directories and perform file transfers using `rsync`. The app also includes the functionality to stop the ongoing file transfer at any point.

## Features

- **Select Source and Destination Paths**: Easily select the directories you want to sync using a graphical file dialog.
- **Rsync Transfer**: Transfers files from the source to the destination using the powerful `rsync` command.
- **Stop Ongoing Transfer**: You can stop an ongoing transfer by clicking the "Stop Rsync" button.

## Prerequisites

- **Go**: Make sure Go is installed. You can download it from [here](https://golang.org/dl/).
- **rsync**: The `rsync` command must be installed on your system.

## Installation

1. Clone this repository or copy the source code.

   ```bash
   git clone https://github.com/rgritt/go-rsync-gui.git
   cd go-rsync-gui
   ```

2. Install the required Go dependencies, including the Fyne UI toolkit:

   ```bash
   go get fyne.io/fyne/v2
   go get
   go mod tidy
   ```

3. Run the app:

   ```bash
   go run main.go
   ```

## Usage

1. **Select Source Directory**: Click on the "Select Source Directory" button and choose the directory from which you want to transfer files.
2. **Select Destination Directory**: Click on the "Select Destination Directory" button and choose the destination directory where the files should be transferred.
3. **Start Rsync**: Once the directories are selected, click the "Start Rsync" button to begin the transfer.
4. **Stop Rsync**: If needed, click the "Stop Rsync" button to halt the ongoing file transfer.

## Stop Transfer

The app includes a "Stop Rsync" button that stops the `rsync` process in progress. This can be useful if you need to cancel the transfer for any reason.

## Code Structure

- **main.go**: The primary application logic, including UI components and rsync process handling.
- **runRsync()**: Function that constructs and runs the `rsync` command.
- **UI Components**: Built using the Fyne UI library, providing buttons and dialogs to control the transfer.

## Troubleshooting

1. **Rsync not found**: Make sure `rsync` is installed and available in your system's `PATH`.
   - On Ubuntu/Debian-based systems: `sudo apt install rsync`
   - On macOS (with Homebrew): `brew install rsync`

2. **Go errors**: Ensure all dependencies are properly installed using `go get` before running the program.

## Contributing

Contributions are welcome! If you have ideas for new features or improvements, feel free to fork the repository and submit a pull request.