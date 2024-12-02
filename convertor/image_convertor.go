package converter

import (
	"fmt"

	"github.com/JoelVCrasta/goskii/generator"
	"github.com/JoelVCrasta/goskii/imageutils"
)

func ImageToASCII(path string) error {
	imageData, err := imageutils.LoadImage(path)
	if err != nil {
		return fmt.Errorf("load error: %v", err)
	}

	width, height, err := imageutils.CalculateNewBounds(imageData.Width, imageData.Height)
	if  err != nil {
		return fmt.Errorf("bounds error: %v", err)
	}

	imageGray := imageutils.Grayscale(imageData.Image)
	resizedImage := imageutils.BilinearResizeGray(imageGray, imageData.Bounds, width, height)

	ascii := generator.GenerateASCII(resizedImage, width, height, generator.CHARSET_1)

	fmt.Println(ascii)

	return nil
}