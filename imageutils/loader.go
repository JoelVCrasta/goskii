package imageutils

import (
	"fmt"
	"image"
	"os"
	"strings"
)

type ImageInfo struct {
	Path string
	Image image.Image
	Bounds image.Rectangle
	Width int
	Height int
	FileName string
}

func LoadImage(path string) (*ImageInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding file: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	fileName := path[strings.LastIndex(path, "/")+1:]

	return &ImageInfo{
        Path:     path,
        Image:    img,
        Width:    width,
        Height:   height,
        Bounds:   bounds,
        FileName: fileName,
    }, nil
}