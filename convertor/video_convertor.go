package convertor

import (
	"fmt"
	"image"
	"strings"
	"sync"

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

	var builder strings.Builder
	asciiFrames := make(chan string, videoData.TotalFrames)
	var wg sync.WaitGroup

	// Progress Bar
	bar := progressbar.NewOptions(
		videoData.TotalFrames,
		progressbar.OptionSetPredictTime(false),	
	)

	for i, frame := range videoData.Video {
		wg.Add(1)

		go func(frame image.Image) {
			defer wg.Done()

			frameGray := utils.Grayscale(frame)
			resizedFrame := utils.ResizeGray(frameGray, width, height)
			asciiFrame := generator.GenerateASCII(resizedFrame, width, height, flags.Charset - 1)

			asciiFrames <- asciiFrame

			bar.Describe(fmt.Sprintf("%d/%d frame(s)", i+1, videoData.TotalFrames))
			_ = bar.Add(1)
		}(frame)
	}

	go func() {
		wg.Wait()
		close(asciiFrames)
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