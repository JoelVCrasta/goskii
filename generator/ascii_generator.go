package generator

import (
	"image"
	"image/color"
	"strings"
)

var charsets = [][]string{
	// ASCII 1 lightest to darkest
    {".", ":", "-", "+", "*", "?", "#", "%", "$", "@"},

	// ASCII 1 darkest to lightest
    {"@", "$", "%", "#", "?", "*", "+", "-", ":", "."},

	// ASCII 2 lightest to darkest
	{"`", ".", ",", "\"", ":", ";","~", "^", "?", "#",},

	// ASCII 2 darkest to lightest
	{"?", "!", "^", "~", ";", ":", "\"", ",", ".", "`"},

	// ASCII 3 mixed
	{"q", "#", "!", "?", "o", "%", "*", "+", "-", "="},

	// ASCII 4 mixed
	{"<", ">", ":", ";", "I", "1", "!", "i", ":", "."},

	// Numeric
	{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},

	// Alphabetic
	{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"},

	// Alphanumeric
	{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",},

	// Code Page 437 lightest to darkest
    {" ", "░", "▒", "▓", "█"},

	// Code Page 437 darkest to lightest
    {"█", "▓", "▒", "░", " "},

	// Unicode Mathemathical Operators
	{"+", "-", "*", "/", "=", "<", ">", "±", "÷", "≈", "≠", "≤", "≥", "∞"},

	// Unicode Arrows
	{"←", "↑", "→", "↓", "↔", "↕", "⇐", "⇑", "⇒", "⇓"},
}


func getASCIIChar(c color.Gray, charset int) string {
	gray := c.Y
	normalized := int(float64(gray) / 255.0 * float64(len(charsets[charset]) - 1))

	return charsets[charset][normalized]
}

func GetCharsets () [][]string {
	return charsets
}

// Generates ASCII art from a grayscale image.
func GenerateASCII(img *image.Gray, width, height int, charset int) string {
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

// Generates ASCII art from a grayscale image with alpha channel.
func GenerateASCIIAlpha(img *image.Gray, alpha [][]uint8, width, height int, charset int) string {
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