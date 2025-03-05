package convertor

import (
	"fmt"
	"image"

	"github.com/JoelVCrasta/goskii/cmd"
	"github.com/JoelVCrasta/goskii/generator"
	"github.com/JoelVCrasta/goskii/utils"
)

// Converts an image to grayscale, resizes it, and generates ASCII art.
func convertImage(imageData *utils.ImageData, width, height, charset int, hasAlpha bool) string {
	var imageGray *image.Gray
	var alpha [][]uint8

	if hasAlpha {
		imageGray, alpha = utils.GrayscaleAlpha(imageData.Image)
		alpha = utils.ResizeAlpha(alpha, imageData.Width, imageData.Height, width, height)
	} else {
		imageGray = utils.Grayscale(imageData.Image)
	}

	resizedImage := utils.ResizeGray(imageGray, width, height)

	if hasAlpha {
		return generator.GenerateASCIIAlpha(resizedImage, alpha, width, height, charset)
	}
	return generator.GenerateASCII(resizedImage, width, height, charset-1)
}


// Converts the image to ASCII by calling the appropriate function based on the image extension.
func ImageToASCII(
	flags cmd.Command,
) error {
	imageData, err := utils.LoadImage(flags.Path)
	if err != nil {
		return fmt.Errorf("load error: %v", err)
	}

	width, height, err := utils.CalculateNewBounds(imageData.Width, imageData.Height, flags.Size)
	if  err != nil {
		return fmt.Errorf("bounds error: %v", err)
	}

	termW, termH, err := utils.GetTerminalSize()
	if err != nil {
		return fmt.Errorf("terminal size error: %v", err)
	}

	shouldPrint := width <= termW && height <= termH
	if !shouldPrint && flags.Output == "" {
		fmt.Println("ASCII art is too large to fit in the terminal. Increase the terminal size or use the -o flag to save to a file.")
	}


	var ascii string
	if imageData.Extension == ".png" {
		ascii = convertImage(imageData, width, height, flags.Charset, true)
	} else {
		ascii = convertImage(imageData, width, height, flags.Charset, false)
	}

	if shouldPrint {
		fmt.Println(ascii)
	}

	if flags.Output != "" {
		err := utils.SaveToTextFile(ascii, flags.Output, imageData.FileName)
		if err != nil {
			return fmt.Errorf("save error: %v", err)
		}
	}

	return nil
}