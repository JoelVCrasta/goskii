package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

func waitKeyPress() {
	prevState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), prevState)

	buf := make([]byte, 1)
	loop:
		for {
			_, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
	
			switch buf[0] {
			case 13: break loop // Enter
			case 27, 81, 113: os.Exit(0) // Escape | Q | q
			}
		}
}

func Render(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	firstLine := strings.SplitN(string(content), "\n", 2)[0]
	lineWidth := len(firstLine)

	termW, _, err := GetTerminalSize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if lineWidth > termW {
		widthDiff := lineWidth - termW
		fmt.Printf("The ASCII is %d characters wider than the terminal. Resize the terminal to fit the ASCII art.\n", widthDiff)
		fmt.Println("Press 'Enter' to continue...")
		waitKeyPress()
	}

	checkVideo := strings.Split(string(content), "\n\n")
	if len(checkVideo) > 1 {
		RenderVideo(string(content))
	} else {
		fmt.Print(string(content))
	}
}

func RenderVideo(ascii string) {
	frames := strings.Split(ascii, "\n\n")
	frameDelay := time.Second / 12
	ClearTerminal()

	for _, frame := range frames {
		fmt.Print("\033[H")
		fmt.Print(frame)
		time.Sleep(frameDelay)
	}

	ClearTerminal()
}