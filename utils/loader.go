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
	"os/exec"
	"path/filepath"
	"regexp"
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

// extractVideoId extracts the video ID from a YouTube URL.
func extractVideoId(url string) (string, error) {
	re := regexp.MustCompile(`(?:=|be\/)([0-9A-Za-z_-]{11})`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("video ID not found in URL: %s", url)
	}

	return matches[1], nil
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
    )

	var metadata struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	if strings.Contains(path, "youtube.com") || strings.Contains(path, "youtu.be") {
		// Handle youtube path
		fetchQuality := "360"

		videoId, err := extractVideoId(path)
		if err != nil {
			return nil, err
		}
		outputTemplate := videoId + "-goskii.%(ext)s"

		cmd := exec.Command(
			"yt-dlp",
			"-f", "bestvideo[height<="+fetchQuality+"][ext=mp4]",
			"--concurrent-fragments", "4",
			"-o", outputTemplate,
			path,
		)
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("error fetching YouTube video: %v", err)
		}

		matches, _ := filepath.Glob(videoId + "-goskii.*")
		if len(matches) == 0 {
			return nil, fmt.Errorf("dowloaded youtube video not found")
		}
		path = matches[0]
    }

	// Handle HTTP/HTTPS, local file and downloaed youtube video
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

		if err != nil {
			writer.CloseWithError(fmt.Errorf("error reading video: %v", err))
		}
	}()

    return &VideoData{
        Path:      path,
        Reader:    reader,
        Width:     width,
        Height:    height,
        FileName:  strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
        Extension: filepath.Ext(path),
    }, nil
}