package generator

import (
	"image"
	"image/color"
	"strings"
)

// ASCII character sets to represent the image.
var (
	// The characters are ordered from darkest to lightest.
	CHARSET_1 = []string{".", ":", "-", "+", "*", "?", "#", "%", "$", "@"}

	// The characters are ordered from lightest to darkest.
	CHARSET_2 = []string{"@", "$", "%", "#", "?", "*", "+", "-", ":", "."}
)

func getASCIIChar(c color.Gray, charset []string ) string {
	gray := c.Y
	normalized := int(float64(gray) / 255.0 * float64(len(charset)-1))

	return CHARSET_1[normalized]
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