package imageutils

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

type ImageData struct {
	Path string
	Image image.Image
	Width int
	Height int
	FileName string
	Extension string
}

func LoadImage(path string) (*ImageData, error) {
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
	extension := strings.ToLower(path[strings.LastIndex(path, ".")+1:])

	return &ImageData{
        Path:     path,
        Image:    img,
        Width:    width,
        Height:   height,
        FileName: fileName,
		Extension: extension,
    }, nil
}