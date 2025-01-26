package utils

import (
	"bytes"
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
	Path      string
	Image     image.Image
	Width     int
	Height    int
	FileName  string
	Extension string
}

type VideoData struct {
	Path      string
	Video  	  []image.Image
	Width	  int
	Height	  int
	FileName  string
	Extension string
	
}

// decodeMJPEGStream reads an MJPEG stream from the reader which is read from ffmpeg output and decodes it into a slice of images.
func decodeMJPEGStream(reader io.Reader) ([]image.Image, error) {
	var frames []image.Image
	frameBuffer := bytes.Buffer{}
	buf := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			frameBuffer.Write(buf[:n])

			data := frameBuffer.Bytes()
			startIdx := bytes.Index(data, []byte{0xFF, 0xD8}) // SOI marker
			endIdx := bytes.Index(data, []byte{0xFF, 0xD9})   // EOI marker

			if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
				jpegData := data[startIdx : endIdx+2]

				// Decode the JPEG frame
				img, err := jpeg.Decode(bytes.NewReader(jpegData))
				if err != nil {
					return nil, fmt.Errorf("failed to decode JPEG frame: %w", err)
				}

				frames = append(frames, img)

				// Remove processed frame data from the buffer
				frameBuffer = *bytes.NewBuffer(data[endIdx+2:])
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading MJPEG stream: %w", err)
		}
	}

	return frames, nil
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

func LoadVideo(path string) (*VideoData, error) {
	reader, writer := io.Pipe()

	var frames []image.Image

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

	frames, err := decodeMJPEGStream(reader)

	if err != nil {
		return nil, fmt.Errorf("error decoding MJPEG stream: %v", err)
	}

	return &VideoData{
		Path: path,
		Video: frames,
		Width: frames[0].Bounds().Dx(),
		Height: frames[0].Bounds().Dy(),
		FileName: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Extension: filepath.Ext(path),
	}, nil

}