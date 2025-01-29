package utils

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type ImageData struct {
	Path      	string
	Image     	image.Image
	Width     	int
	Height    	int
	FileName  	string
	Extension 	string
}

type VideoData struct {
	Path      	string
	Reader 		*io.PipeReader
	Width	  	int
	Height	  	int
	FileName  	string
	Extension 	string
}

func LoadImage(path string) (*ImageData, error) {
	var reader io.Reader

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
	} else {
		// Handle local file input
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		reader = file
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
		FileName:  strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Extension: filepath.Ext(path),
	}, nil
}

// TODO: Support for URL videos
func LoadVideo(path string) (*VideoData, error) {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()

		err := ffmpeg.Input(path).Output(
			"pipe:1", ffmpeg.KwArgs{
				"format": "image2pipe",
				"vcodec": "mjpeg",
				"r": "12",
			},	
		).WithOutput(writer).Run()

		if err != nil {
			fmt.Printf("FFmpeg error: %v\n", err)
		}
	}()

	img, err := jpeg.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	return &VideoData{
		Path: path,
		Reader: reader,
		Width: img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		FileName: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Extension: filepath.Ext(path),
	}, nil
}