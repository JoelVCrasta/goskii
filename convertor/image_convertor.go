package convertor

import (
	"fmt"

	"github.com/JoelVCrasta/goskii/generator"
	"github.com/JoelVCrasta/goskii/imageutils"
)

// Converts the Non Alpha image to ASCII.
func imageRGB(imageData *imageutils.ImageData, width int, height int) {
	imageGray := imageutils.Grayscale(imageData.Image)
	resizedImage := imageutils.BilinearResizeGray(imageGray, width, height)

	ascii := generator.GenerateASCII(resizedImage, width, height, generator.CHARSET_1)
	fmt.Println(ascii)
}

// Converts the Alpha image to ASCII. This is for PNG images.
func imageRGBA(imageData *imageutils.ImageData, width int, height int) {
	imageGray, alpha := imageutils.GrayscaleAlpha(imageData.Image)
	resizedImage := imageutils.BilinearResizeGray(imageGray, width, height)
	alpha = imageutils.ResizeAlpha(alpha, imageData.Width, imageData.Height, width, height)

	ascii := generator.GenerateASCIIAlpha(resizedImage, alpha, width, height, generator.CHARSET_1)
	fmt.Println(ascii)
}

// Converts the image to ASCII by calling the appropriate function based on the image extension.
func ImageToASCII(path string) error {
	imageData, err := imageutils.LoadImage(path)
	if err != nil {
		return fmt.Errorf("load error: %v", err)
	}

	width, height, err := imageutils.CalculateNewBounds(imageData.Width, imageData.Height)
	if  err != nil {
		return fmt.Errorf("bounds error: %v", err)
	}

	if imageData.Extension == "png" {
		imageRGBA(imageData, width, height)
	} else {
		imageRGB(imageData, width, height)
	}

	return nil
}