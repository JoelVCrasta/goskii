package utils

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkdai/youtube/v2"
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

// Returns the format with the specified quality
func findFormat(formats youtube.FormatList, quality string) *youtube.Format {
	for i := range formats {
		if strings.Contains(formats[i].QualityLabel, quality) {
			return &formats[i]
		}
	}
		
	if len(formats) > 0 {
		return &formats[0]
	}

	return nil
}

// LoadImage loads an image from the specified path (local or http) and returns an ImageData struct containing the image and metadata.
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

// LoadVideo loads a video from the specified path (local, http or youtube) and returns a VideoData struct containing the video stream and metadata.
func LoadVideo(path string) (*VideoData, error) {
    var (
        reader, writer  = io.Pipe()
        width, height   = 0, 0
        errChan         = make(chan error, 1)
    )

	var metadata struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	if strings.Contains(path, "youtube.com") || strings.Contains(path, "youtu.be") {
		// TODO: Implement YouTube video support
		client := youtube.Client{}
		fetchQuality := "360p"

		video, err := client.GetVideo(path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch YouTube video: %v", err)
		}

		fmt.Printf("Downloading YouTube video: %s\n", video.Title)

		format := findFormat(video.Formats, fetchQuality)
		if format == nil {
			return nil, fmt.Errorf("no format found for quality '%s'", fetchQuality)
		}
		height, width = format.Height, format.Width

		stream, _, err := client.GetStream(video, format)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch video stream: %v", err)
		}

		go func() {
			fmt.Println("Starting FFmpeg")
			defer writer.Close()
			err := ffmpeg.Input("pipe:0").Output(
				"pipe:1", ffmpeg.KwArgs{
					"format": "image2pipe",
					"vcodec": "mjpeg",
					"r":      "12",
				},
			).WithInput(stream).WithOutput(writer).ErrorToStdOut().Run()
			if err != nil {
				errChan <- fmt.Errorf("ffmpeg error: %v", err)
				return
			}
			errChan <- nil
			fmt.Println("FFmpeg finished")
		}()
    } else {
		// Handle HTTP/HTTPS input (non-YouTube) and local file input
		probeResult, err := ffmpeg.Probe(path)
		if err != nil {
			return nil, fmt.Errorf("error probing video: %v", err)
		}
		if err := json.Unmarshal([]byte(probeResult), &metadata); err != nil {
			return nil, fmt.Errorf("error parsing video metadata: %v", err)
		}
		if len(metadata.Streams) == 0 {
			return nil, fmt.Errorf("no video streams found")
		}
		width, height = metadata.Streams[0].Width, metadata.Streams[0].Height

		go func() {
			defer writer.Close()
			err := ffmpeg.Input(path).Output(
				"pipe:1", ffmpeg.KwArgs{
					"format": "image2pipe",
					"vcodec": "mjpeg",
					"r":      "12",
				},
			).WithOutput(writer).Silent(true).Run()

			errChan <- err
		}()
	}

	if err := <-errChan; err != nil {
		return nil, err
	}

    return &VideoData{
        Path:      path,
        Reader:    reader,
        Width:     width,
        Height:    height,
        FileName:  strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
        Extension: filepath.Ext(path),
    }, nil
}