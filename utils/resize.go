package utils

import (
	"image"
	"image/color"
	"math"
)

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Bilinear interpolation function
func bilinearInterpolate(x, y float64, x1, y1, x2, y2 float64, q11, q12, q21, q22 uint8) uint8 {
	r1 := ((x2 - x) / (x2 - x1)) * float64(q11) + ((x - x1) / (x2 - x1)) * float64(q21)
	r2 := ((x2 - x) / (x2 - x1)) * float64(q12) + ((x - x1) / (x2 - x1)) * float64(q22)
	result := ((y2 - y) / (y2 - y1)) * r1 + ((y - y1) / (y2 - y1)) * r2

	return uint8(math.Round(result))
}

// ResizeGray resizes a grayscale image using bilinear interpolation
func ResizeGray(img *image.Gray, newWidth, newHeight int) *image.Gray {
	origWidth := img.Bounds().Dx()
	origHeight := img.Bounds().Dy()

	resizedImage := image.NewGray(image.Rect(0, 0, newWidth, newHeight))
	xRatio := float64(origWidth) / float64(newWidth)
	yRatio := float64(origHeight) / float64(newHeight)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			origX := float64(x) * xRatio
			origY := float64(y) * yRatio

			// Find the four surrounding pixels
			x1 := int(origX)
			y1 := int(origY)
			x2 := min(x1+1, origWidth-1)
			y2 := min(y1+1, origHeight-1)

			// Get pixel values
			q11 := img.GrayAt(x1, y1).Y
			q21 := img.GrayAt(x2, y1).Y
			q12 := img.GrayAt(x1, y2).Y
			q22 := img.GrayAt(x2, y2).Y

			// Perform bilinear interpolation
			lum := bilinearInterpolate(origX, origY, float64(x1), float64(y1), float64(x2), float64(y2), q11, q12, q21, q22)
			resizedImage.SetGray(x, y, color.Gray{Y: lum})
		}
	}

	return resizedImage
}

// ResizeAlpha resizes a 2D alpha channel array using bilinear interpolation
func ResizeAlpha(alpha [][]uint8, origWidth, origHeight, newWidth, newHeight int) [][]uint8 {
	resizedAlpha := make([][]uint8, newHeight)
	for i := range resizedAlpha {
		resizedAlpha[i] = make([]uint8, newWidth)
	}

	xRatio := float64(origWidth) / float64(newWidth)
	yRatio := float64(origHeight) / float64(newHeight)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			origX := float64(x) * xRatio
			origY := float64(y) * yRatio

			// Find the four surrounding pixels
			x1 := int(origX)
			y1 := int(origY)
			x2 := min(x1+1, origWidth-1)
			y2 := min(y1+1, origHeight-1)

			// Get intensity values
			q11 := alpha[y1][x1]
			q21 := alpha[y1][x2]
			q12 := alpha[y2][x1]
			q22 := alpha[y2][x2]

			// Perform bilinear interpolation
			resizedAlpha[y][x] = bilinearInterpolate(origX, origY, float64(x1), float64(y1), float64(x2), float64(y2), q11, q12, q21, q22)
		}
	}

	return resizedAlpha
}
