package convertor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/JoelVCrasta/goskii/cmd"
	"github.com/JoelVCrasta/goskii/generator"
	"github.com/JoelVCrasta/goskii/utils"
)

// processFrames processes a batch of frames concurrently and appends the ASCII representation to the builder.
func processFrames(frames []image.Image, builder *strings.Builder, charset, width, height int, frameCount *int32) {
	var (
		wg 			sync.WaitGroup
		asciiFrames = make([]string, len(frames))	
	)

	for idx, frame := range frames {
		wg.Add(1)
		go func(i int, f image.Image) {
			defer wg.Done()

			grayFrame := utils.Grayscale(f)
			resizedFrame := utils.ResizeGray(grayFrame, width, height)
			asciiFrames[i] = generator.GenerateASCII(resizedFrame, width, height, charset-1)

			atomic.AddInt32(frameCount, 1)
		}(idx, frame)
	}

	wg.Wait()

	for _, asciiFrame := range asciiFrames {
		builder.WriteString(asciiFrame)
		builder.WriteString("\n\n")
	}
}

/* 
	decodeAndProcessStream extracts frames from an MJPEG stream, converts them to ASCII, and returns the final ASCII representation.

	1) Read the stream from io.PipeReader in chunks and accumulate data in a buffer.

	2) Detect JPEG frames using SOI (0xFFD8) and EOI (0xFFD9) markers.

	3) Decode and append frames to the frames slice and process them in batches of 12.

	4) If slice length is 12, process the frames concurrently and reset the slice.

	5) Process the remaining frames and return the final ASCII representation.
*/
func decodeAndProcessStream(videoData *utils.VideoData, charset, width, height int) (string, error) { 
	var (
		frames 		= make([]image.Image, 0, 16) // Slice to store 16 frames
		builder	 	strings.Builder
		frameCount 	int32
		frameBuffer bytes.Buffer
		buf 		= make([]byte, 1024)
	)

	for {
		n, err := videoData.Reader.Read(buf)
		if n > 0 {
			frameBuffer.Write(buf[:n])

			data := frameBuffer.Bytes()
			startIdx := bytes.Index(data, []byte{0xFF, 0xD8}) // SOI marker
			endIdx := bytes.Index(data, []byte{0xFF, 0xD9})   // EOI marker

			if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
				jpegData := data[startIdx : endIdx+2]

				// Decode the JPEG frame and append it to the frames slice
				frame, err := jpeg.Decode(bytes.NewReader(jpegData))
				if err != nil {
					return "", fmt.Errorf("failed to decode JPEG frame: %w", err)
				}
				frames = append(frames, frame)

				// Remove processed frame data from the buffer and keep the unprocessed data
				frameBuffer.Next(endIdx + 2)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading MJPEG stream: %w", err)
		}

		if len(frames) == 16 {
			processFrames(frames, &builder, charset, width, height, &frameCount)
			frames = frames[:0]
		}
	}

	if len(frames) > 0 {
		processFrames(frames, &builder, charset, width, height, &frameCount)
	}

	return builder.String(), nil
}


// VideoToASCII converts a video to ASCII art.
func VideoToASCII(flags cmd.Command) error {
	videoData, err := utils.LoadVideo(flags.Path)
	if err != nil {
		return fmt.Errorf("load error: %v", err)
	}
	defer videoData.Reader.Close()

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

	ascii, err := decodeAndProcessStream(videoData, flags.Charset, width, height)
	if err != nil {
		return fmt.Errorf("error processing stream: %v", err)
	}

	if !shouldPrint && ascii != "" {
		utils.RenderVideo(ascii, flags.Fps)
	}

	if flags.Output != "" {
		utils.SaveToTextFile(ascii, flags.Output, videoData.FileName)
	}

	return nil
}
