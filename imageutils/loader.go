package imageutils

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type ImageData struct {
	Path      string
	Image     image.Image
	Width     int
	Height    int
	FileName  string
	Extension string
}

func LoadImage(path string) (*ImageData, error) {
	var reader io.Reader
	var extension string
	var fileName string

	if strings.HasPrefix(path, "http") {
		// Handle URL input
		res, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("error fetching URL: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error fetching URL, status code: %d", res.StatusCode)
		}

		reader = res.Body
		extension = filepath.Ext(path)
		fileName = strings.TrimSuffix(filepath.Base(path), extension)
	} else {
		// Handle local file input
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		reader = file
		extension = filepath.Ext(path)
		fileName = strings.TrimSuffix(filepath.Base(path), extension)
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	return &ImageData{
		Path:      path,
		Image:     img,
		Width:     img.Bounds().Dx(),
		Height:    img.Bounds().Dy(),
		FileName:  fileName,
		Extension: extension,
	}, nil
}
