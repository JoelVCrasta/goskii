package imageutils

import (
	"image"
	"image/color"
)

// Grayscale converts an image to grayscale using the formula:
// gray = 0.299*r + 0.587*g + 0.114*b
func Grayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, _ := originalColor.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			grayImg.Set(x, y, pixel)
		}
	}

	return grayImg
}	