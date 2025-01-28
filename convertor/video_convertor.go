package convertor

import (
	"fmt"
	"strings"

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
	totalFrames := len(videoData.Video)

	bar := progressbar.NewOptions(len(videoData.Video),
		progressbar.OptionSetPredictTime(false),	
	)

	for i, frame := range videoData.Video {
		frameGray := utils.Grayscale(frame)
		resizedFrame := utils.ResizeGray(frameGray, width, height)

		asciiFrame := generator.GenerateASCII(resizedFrame, width, height, flags.Charset - 1)
		builder.WriteString(asciiFrame)
		builder.WriteString("\n\n")

		bar.Describe(fmt.Sprintf("%d/%d frame(s)", i+1, totalFrames))
		bar.Add(1)
	}
	bar.Clear()

	if flags.Output != "" {
		utils.SaveToTextFile(builder.String(), flags.Output, videoData.FileName)
	}

	return nil
}