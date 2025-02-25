package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JoelVCrasta/goskii/generator"
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
)

const (
	DefaultSize     = 0
    MinSize         = 0
    MaxSize         = 500
    DefaultCharset  = 1
    MinCharset      = 1
    MaxCharset      = 13
	MinFps			= 1
	MaxFps			= 24
	DefaultFps		= 12
	Version 		= "2.0"
)

type Command struct {
	Path   			string
	Output 			string
	Render  		string
	Size  			int
	Charset 		int
	Fps 			int
}
var cmdFlags Command

var rootCmd = &cobra.Command{
	Use:  "goskii",
	Short: "goskii is a CLI tool to convert images to ASCII art.",

	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			fmt.Println("Usage: goskii [flags]")
			fmt.Println("Type --help or -h to see a list of all options.")   
			os.Exit(0)
		}

		if cmd.Flags().Changed("showset") {
			showShowset()
			os.Exit(0)
		}

		if cmd.Flags().Changed("version") {
			fmt.Printf("goskii version %s\n", Version)
			os.Exit(0)
		}

		if !checkFilePath(cmd, &cmdFlags.Path) && !checkRender(cmd, &cmdFlags.Render) {
			os.Exit(1)
		}

		if !checkOutputPath(cmd, &cmdFlags.Output) {
			os.Exit(1)
		}

		if !checkSize(cmd, &cmdFlags.Size) {
			os.Exit(1)
		}

		if !checkCharset(cmd, &cmdFlags.Charset) {
			os.Exit(1)
		}

		if !checkFps(cmd, &cmdFlags.Fps, &cmdFlags.Path, &cmdFlags.Render) {
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&cmdFlags.Path, "path", "p", "","Path to the file. (Required)")
    rootCmd.Flags().StringVarP(&cmdFlags.Output, "output", "o", "", "Output folder path. Default is current directory.")
	rootCmd.Flags().StringVarP(&cmdFlags.Render, "render", "r", "", "Render the contents of the ASCII art file.")
    rootCmd.Flags().IntVarP(&cmdFlags.Size, "width", "w", DefaultSize, fmt.Sprintf("Width of the ASCII art (%d - %d). Default adjusts to terminal size.", MinSize, MaxSize))
    rootCmd.Flags().IntVarP(&cmdFlags.Charset, "charset", "c", DefaultCharset, fmt.Sprintf("Character set to use (%d - %d).", MinCharset, MaxCharset))
	rootCmd.Flags().IntVarP(&cmdFlags.Fps, "fps", "f", 12, fmt.Sprintf("Video FPS (%d - %d). Default is %d.", MinFps, MaxFps, DefaultFps))
    rootCmd.Flags().BoolP("showset", "s", false, "Display all character sets.")
	rootCmd.Flags().BoolP("version", "v", false, "Verion of goskii.")
	rootCmd.MarkPersistentFlagRequired("path")
	
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrln(err)
		os.Exit(1)
	}
} 

// Returns all the command line arguments.
func GetCommands() Command {
	return cmdFlags
}

// Returns the file type (image or video)
func GetFileType() int {
	return checkExtension(cmdFlags.Path)
}

// Checks whether the file extension is supported.
func checkExtension(path string) int {
	ext := strings.ToLower(filepath.Ext(path))
	youtubeRegex := regexp.MustCompile(`https?:\/\/(www\.)?(youtube\.com|youtu\.be)\/`)
	if youtubeRegex.MatchString(path) {
		return 3
	}

	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".tiff", ".bmp", ".gif":
		return 0
	case ".mp4", ".avi", ".mov", ".mkv", ".flv", ".webm", ".mpeg":
		return 1
	case ".txt":
		return 2
	default:
		return -1
	}
}

// Checks whether the file path is valid and the file extension is supported.
func checkFilePath(cmd *cobra.Command, path *string) bool {
	if *path == "" {
		return false
	}

	// Check if the path is a valid URL
	parsedURL, err := url.ParseRequestURI(*path)
	if err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		// Handle HTTP/HTTPS URLs
		if strings.Contains(parsedURL.Host, "youtube.com") || strings.Contains(parsedURL.Host, "youtu.be") {
			// Check if it's a valid YouTube video
			client := youtube.Client{}
			_, err := client.GetVideo(*path)
			if err != nil {
				cmd.PrintErrf("Error fetching the YouTube video: %v\n", err)
				return false
			}
			return true
		}

		// Handle non-YouTube URLs (e.g., images or videos)
		res, err := http.Head(*path)
		if err != nil {
			cmd.PrintErrf("Error fetching the URL \"%s\": %v\n", *path, err)
			return false
		}
		defer res.Body.Close()

		contentType := res.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "video/") {
			cmd.PrintErrf("The URL \"%s\" does not point to an image or a video.\n", *path)
			return false
		}

		return true
	}

	// Handle local file paths
	if _, err := os.Stat(*path); os.IsNotExist(err) {
		cmd.PrintErrf("The path or file \"%s\" does not exist or is not valid.\n", *path)
		return false
	}

	if checkExtension(*path) == -1 {
		cmd.PrintErrf("The file extension is not supported.\n")
		return false
	}

	return true
}

// Checks whether the output path is valid and a directory.
func checkOutputPath(cmd *cobra.Command, path *string) bool {
	if *path == "" {
        return true
    }

    if *path == "." {
        cwd, err := os.Getwd()
        if err != nil {
            cmd.PrintErrf("Error getting current working directory: %v\n", err)
            return false
        }
        *path = cwd
        return true
    }

    info, err := os.Stat(*path)
    if os.IsNotExist(err) {
        cmd.PrintErrf("The output folder \"%s\" does not exist.\n", *path)
        return false
    } else if err != nil {
        cmd.PrintErrf("Error: %v\n", err)
        return false
    } else if !info.IsDir() {
        cmd.PrintErrf("The output path \"%s\" is not a directory.\n", *path)
        return false
    }


    return true
}

func checkRender(cmd *cobra.Command, path *string) bool {
	if *path == "" {
		return false
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		cmd.PrintErrf("The path or file \"%s\" does not exist or not valid.\n", *path)
		return false
	}

	if filepath.Ext(*path) != ".txt" {
		cmd.PrintErrf("The file extension is not supported.\n")
		return false
	}

	return true
}

// Checks whether the size is between 1 and 500.
func checkSize(cmd *cobra.Command, size *int) bool {
	if *size < MinSize || *size > MaxSize {
		cmd.PrintErrf("The size should be between %d and %d.", MinSize+1, MaxSize)
		return false
	}

	return true
}

// Checks whether the charset is between 1 and 10.
func checkCharset(cmd *cobra.Command, charset *int) bool {
	if *charset < MinCharset || *charset > MaxCharset {
		cmd.PrintErrf("The charset should be between %d and %d.\n", MinCharset, MaxCharset)
		return false
	}

	return true
}

func checkFps(cmd *cobra.Command, fps *int, path *string, render *string) bool {
	if *path == "" && *render == "" {
		cmd.PrintErrf("No video path provided.\n")
		return false
	}

	if checkExtension(*path) != 1 && checkExtension(*render) != 2 && strings.Contains(*path, "youtube.com") {
		cmd.PrintErrf("Not a video file.\n")
		return false
	}

	if *fps < MinFps || *fps > MaxFps {
		cmd.PrintErrf("The FPS should be between %d and %d.\n", MinFps, MaxFps)
		return false
	}

	return true
}

// Displays all the character sets.
func showShowset() {
	charsetsDesc := map[int]string{
		1: "ASCII 1 Lightest to darkest characters. Better for terminal. (Default)",
		2: "ASCII 1 Darkest to lightest characters. Better for light background.",
		3: "ASCII 2 Lightest to darkest characters.",
		4: "ASCII 2 Darkest to lightest characters.",
		5: "ASCII 3 Mixed.",
		6: "ASCII 4 Mixed.",
		7: "Numeric characters.",
		8: "Alphabetic characters.",
		9: "Alphanumeric characters.",
		10: "Lightest to darkest block shades (Code Page 437).",
		11: "Darkest to lightest block shades (Code Page 437).",
		12: "Mathematical Operators (Unicode).",
		13: "Arrows (Unicode).",
	}

	charsets := generator.GetCharsets()
	for i, set := range charsets {
		setStr := strings.Join(set, " ")
		fmt.Printf("%s\n", charsetsDesc[i+1])
		fmt.Printf("%d) %s\n\n", i+1, setStr)
	}

	fmt.Println("Note: The Unicode characters may not work in all terminals.")
}