package imageutils

import (
	"fmt"
	"os"
)

func Render(ascii string) {
	writer := os.Stdout

	_, err := fmt.Fprintln(writer, ascii)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}