package convertor

import (
	"fmt"

	"github.com/JoelVCrasta/goskii/cmd"
	"github.com/JoelVCrasta/goskii/generator"
	"github.com/JoelVCrasta/goskii/imageutils"
)

// Converts the Non Alpha image to ASCII.
func imageRGB(imageData *imageutils.ImageData, width int, height int) string {
	imageGray := imageutils.Grayscale(imageData.Image)
	resizedImage := imageutils.BilinearResizeGray(imageGray, width, height)

	ascii := generator.GenerateASCII(resizedImage, width, height, generator.CHARSET_1)
	fmt.Println(ascii)

	return ascii
}

// Converts the Alpha image to ASCII. This is for PNG images.
func imageRGBA(imageData *imageutils.ImageData, width int, height int) string {
	imageGray, alpha := imageutils.GrayscaleAlpha(imageData.Image)
	resizedImage := imageutils.BilinearResizeGray(imageGray, width, height)
	alpha = imageutils.ResizeAlpha(alpha, imageData.Width, imageData.Height, width, height)

	ascii := generator.GenerateASCIIAlpha(resizedImage, alpha, width, height, generator.CHARSET_1)
	fmt.Println(ascii)

	return ascii
}

// Converts the image to ASCII by calling the appropriate function based on the image extension.
func ImageToASCII(
	flags cmd.Command,
) error {
	imageData, err := imageutils.LoadImage(flags.Path)
	if err != nil {
		return fmt.Errorf("load error: %v", err)
	}

	width, height, err := imageutils.CalculateNewBounds(imageData.Width, imageData.Height, flags.Size)
	if  err != nil {
		return fmt.Errorf("bounds error: %v", err)
	}

	var ascii string
	if imageData.Extension == "png" {
		ascii = imageRGBA(imageData, width, height)
	} else {
		ascii = imageRGB(imageData, width, height)
	}

	if flags.Output != "" {
		imageutils.SaveToTextFile(ascii, flags.Output, imageData.FileName)
	}

	return nil
}