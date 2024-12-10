package imageutils

import (
	"image"
	"image/color"
)

// Grayscale converts an image to grayscale.
func Grayscale(img image.Image) *image.Gray {
    bounds := img.Bounds()
    grayImg := image.NewGray(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            originalColor := img.At(x, y)
            r, g, b, _ := originalColor.RGBA()
            lum := 0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)
            pixel := color.Gray{uint8(lum / 256)}
            grayImg.Set(x, y, pixel)
        }
    }

    return grayImg
}

// GrayscaleAlpha converts an image to grayscale and returns the alpha values.
func GrayscaleAlpha(img image.Image) (*image.Gray, [][]uint8) {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	alpha := make([][]uint8, bounds.Dy())
	for i := range alpha {
		alpha[i] = make([]uint8, bounds.Dx())
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, a := originalColor.RGBA()
			lum := 0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)
			pixel := color.Gray{uint8(lum / 256)}
			grayImg.Set(x, y, pixel)
			alpha[y][x] = uint8(a / 256)
		}
	}

	return grayImg, alpha
}	