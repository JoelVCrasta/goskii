package imageutils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

// GetTerminalSize returns the width and height of the terminal
func GetTerminalSize() (int, int, error) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return 0, 0, fmt.Errorf("not a terminal")
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, fmt.Errorf("error getting terminal size: %v", err)
	}

	return width, height, nil
}


// First, it gets the terminal size
// Then, it calculates the scaling factor for the image based on the terminal size and the image size
// Finally, it calculates the new width and height of the image based on the scaling factor
func CalculateNewBounds(width, height int) (int, int, error) {
	terminalWidth, terminalHeight, err := GetTerminalSize()
	if err != nil {
		return 0, 0, err
	}

	heightScale := 2.0 // default scale to compensate character height of the terminal
	terminalRatio := float64(terminalWidth) / float64(terminalHeight)
	imageRatio := float64(width) / float64(height)

	var scalingFactor float64
	if terminalRatio > imageRatio {
		scalingFactor = float64(terminalHeight) / float64(height)
	} else {
		scalingFactor = float64(terminalWidth) / float64(width)
	}

	newWidth := int(float64(width) * scalingFactor * heightScale)
	newHeight := int(float64(height) * scalingFactor)

	return newWidth, newHeight, nil
}


func ClearTerminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

