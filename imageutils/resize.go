package imageutils

import (
	"image"
	"image/color"
)

// Resize resizes an image to the specified width and height.
func Resize(img image.Image, bounds image.Rectangle, newWidth, newHeight int) (*image.RGBA, error) {
	width, height, err := CalculateNewBounds(bounds.Dx(), bounds.Dy())
	if  err != nil {
		return nil, err
	}

	resizedImage := BilinearResize(img, bounds, width, height)
	return resizedImage, nil
}

// ResizeGray resizes a grayscale image to the specified width and height.
// Faster as it only uses the gray channel of the image.
func ResizeGray(img *image.Gray, bounds image.Rectangle, newWidth, newHeight int) (*image.Gray, error) {
	width, height, err := CalculateNewBounds(bounds.Dx(), bounds.Dy())
	if  err != nil {
		return nil, err
	}

	resizedImage := BilinearResizeGray(img, bounds, width, height)
	return resizedImage, nil
}


func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


func bilinearInterpolate(q11, q21, q12, q22 uint32, dx, dy float64) float64 {
	top := (1-dx)*float64(q11) + dx*float64(q21)
	bottom := (1-dx)*float64(q12) + dx*float64(q22)
	return (1-dy)*top + dy*bottom
}


func BilinearResize(img image.Image, bounds image.Rectangle, newWidth, newHeight int) *image.RGBA {
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	resizedImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	xRatio := float64(origWidth) / float64(newWidth)
	yRatio := float64(origHeight) / float64(newHeight)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			origX := float64(x) * xRatio
			origY := float64(y) * yRatio

			// Find the four closest pixels
			x1 := int(origX)
			y1 := int(origY)
			x2 := min(x1+1, origWidth-1)
			y2 := min(y1+1, origHeight-1)

			// Get the colors of the four surrounding pixels
			c1 := img.At(x1, y1)
			c2 := img.At(x2, y1)
			c3 := img.At(x1, y2)
			c4 := img.At(x2, y2)

			// Interpolate between the pixels
			r1, g1, b1, _ := c1.RGBA()
			r2, g2, b2, _ := c2.RGBA()
			r3, g3, b3, _ := c3.RGBA()
			r4, g4, b4, _ := c4.RGBA()

			// Calculate interpolation weights
			dx := origX - float64(x1)
			dy := origY - float64(y1)

			// Perform bilinear interpolation
			r := bilinearInterpolate(r1, r2, r3, r4, dx, dy)
			g := bilinearInterpolate(g1, g2, g3, g4, dx, dy)
			b := bilinearInterpolate(b1, b2, b3, b4, dx, dy)

			newColor := color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 255}
			resizedImage.Set(x, y, newColor)
		}
	}

	return resizedImage
}


func BilinearResizeGray(img *image.Gray, bounds image.Rectangle, newWidth, newHeight int) *image.Gray {
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	resizedImage := image.NewGray(image.Rect(0, 0, newWidth, newHeight))
	xRatio := float64(origWidth) / float64(newWidth)
	yRatio := float64(origHeight) / float64(newHeight)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			origX := float64(x) * xRatio
			origY := float64(y) * yRatio

			// Find the four closest pixels
			x1 := int(origX)
			y1 := int(origY)
			x2 := min(x1+1, origWidth-1)
			y2 := min(y1+1, origHeight-1)

			// Get the intensity values of the four surrounding pixels
			c1 := img.GrayAt(x1, y1).Y
			c2 := img.GrayAt(x2, y1).Y
			c3 := img.GrayAt(x1, y2).Y
			c4 := img.GrayAt(x2, y2).Y

			// Perform bilinear interpolation
			lum := (float64(c1)*(1-(origX-float64(x1))) + float64(c2)*(origX-float64(x1))) * (1-(origY-float64(y1))) +
				(float64(c3)*(1-(origX-float64(x1))) + float64(c4)*(origX-float64(x1))) * (origY-float64(y1))

			resizedImage.SetGray(x, y, color.Gray{uint8(lum)})
		}
	}

	return resizedImage
}