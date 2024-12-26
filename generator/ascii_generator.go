package generator

import (
	"image"
	"image/color"
	"strings"
)

var (
	// The characters are ordered from lightest to darkest.
	CHARSET_1 = []string{".", ":", "-", "+", "*", "?", "#", "%", "$", "@"}

	// The characters are ordered from darkest to lighest.
	CHARSET_2 = []string{"@", "$", "%", "#", "?", "*", "+", "-", ":", "."}

	// Block shades ordered from lightest to darkest.
	CHARSET_3 = []string{" ", "░", "▒", "▓", "█"}

	// Block shades ordered from darkest to lightest.
	CHARSET_4 = []string{"█", "▓", "▒", "░", " "}
)
var Charsets = [][]string{CHARSET_1, CHARSET_2, CHARSET_3, CHARSET_4}


func getASCIIChar(c color.Gray, charset []string) string {
	gray := c.Y
	normalized := int(float64(gray) / 255.0 * float64(len(charset)-1))

	return charset[normalized]
}


func GenerateASCII(img *image.Gray, width, height int, charset []string) string {
	var builder strings.Builder

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.GrayAt(x, y)
			builder.WriteString(getASCIIChar(pixel, charset))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func GenerateASCIIAlpha(img *image.Gray, alpha [][]uint8, width, height int, charset []string) string {
	var builder strings.Builder

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.GrayAt(x, y)
			alphaValue := alpha[y][x]
			if alphaValue == 0 {
				builder.WriteString(" ")
			} else {
				builder.WriteString(getASCIIChar(pixel, charset))
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}