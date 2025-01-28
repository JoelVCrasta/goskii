package convertor

import (
	"fmt"
	"image"
	"strings"
	"sync"
	"unsafe"

	"github.com/JoelVCrasta/goskii/cmd"
	"github.com/JoelVCrasta/goskii/generator"
	"github.com/JoelVCrasta/goskii/utils"
	"github.com/schollz/progressbar/v3"
)

// Converts the video frames into ASCII.
func VideoToASCII(
	flags cmd.Command,
) error {
	videoData, err := utils.LoadVideo(flags.Path)
	if err != nil {
		return fmt.Errorf("load error: %v", err)
	}

	width, height, err := utils.CalculateNewBounds(videoData.Width, videoData.Height, flags.Size)
	if  err != nil {
		return fmt.Errorf("bounds error: %v", err)
	}

	termW, termH, err := utils.GetTerminalSize()
	if err != nil {
		return fmt.Errorf("terminal size error: %v", err)
	}

	shouldPrint := width > termW || height > termH
	if shouldPrint && flags.Output == "" {
		fmt.Println("ascii art is too large to fit in the terminal, increase the terminal size or use -o flag to save to a file")
	}

	totalSize := 0
    for _, frame := range videoData.Video {
        totalSize += int(unsafe.Sizeof(frame)) + frame.Bounds().Dx()*frame.Bounds().Dy()*4 // Assuming 4 bytes per pixel (RGBA)
    }
    sizeInMB := float64(totalSize) / (1024 * 1024)
    fmt.Printf("Total size of video frames: %d bytes (%.2f MB)\n", totalSize, sizeInMB)

	var builder strings.Builder
	asciiFrames := make(chan string, videoData.TotalFrames)
	progressUpdates := make(chan int, videoData.TotalFrames)
	var wg sync.WaitGroup

	// Progress Bar
	bar := progressbar.NewOptions(
		videoData.TotalFrames,
		progressbar.OptionSetPredictTime(false),	
	)

	go func() {
		completed := 0
		for range progressUpdates {
			completed++
			bar.Describe(fmt.Sprintf("%d/%d frame(s)", completed, videoData.TotalFrames))
			_ = bar.Add(1)
		}
	}()

	for _, frame := range videoData.Video {
		wg.Add(1)

		go func(frame image.Image) {
			defer wg.Done()

			frameGray := utils.Grayscale(frame)
			resizedFrame := utils.ResizeGray(frameGray, width, height)
			asciiFrame := generator.GenerateASCII(resizedFrame, width, height, flags.Charset - 1)

			asciiFrames <- asciiFrame
			progressUpdates <- 1
			//time.Sleep(3000 * time.Millisecond)
		}(frame)
	}

	go func() {
		wg.Wait()
		close(asciiFrames)
		close(progressUpdates)
	}()

	for asciiFrame := range asciiFrames {
		builder.WriteString(asciiFrame)
		builder.WriteString("\n\n")
	}
	bar.Clear()

	if flags.Output != "" {
		utils.SaveToTextFile(builder.String(), flags.Output, videoData.FileName)
	}

	return nil
}