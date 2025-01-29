package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

// Returns the width and height of the terminal.
func GetTerminalSize() (int, int, error) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return 0, 0, fmt.Errorf("not a terminal")
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, fmt.Errorf("%v", err)
	}

	return width, height, nil
}


// Calculates the new width and height of the image based on the terminal size or the size flag.
func CalculateNewBounds(width, height, size int) (int, int, error) {
	terminalWidth, terminalHeight, err := GetTerminalSize()
	if err != nil {
		return 0, 0, err
	}

	heightScale := 2.0 // default scale to compensate character height of the terminal
	var newWidth, newHeight int
	if (size == 0) {
		terminalRatio := float64(terminalWidth) / float64(terminalHeight)
		imageRatio := float64(width) / float64(height)
		
		var scalingFactor float64
		if terminalRatio > imageRatio {
			scalingFactor = float64(terminalHeight) / float64(height) 
		} else {
			scalingFactor = float64(terminalWidth) / float64(width)
		}
		
		newWidth := int(float64(width) * scalingFactor * heightScale)
		newHeight := int(float64(height) * scalingFactor) - 1

		return newWidth, newHeight, nil
	} else {
		newWidth = size
		newHeight = int(float64(height) * float64(newWidth) / (float64(width) * heightScale))

		return newWidth, newHeight, nil
	}
}

// Clears the terminal screen based on the OS.
func ClearTerminal() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}
