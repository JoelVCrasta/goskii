package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	pFlag string
	oFlag string

)

var rootCmd = &cobra.Command{
	Use:  "goskii",
	Short: "goskii is a CLI tool to convert images to ASCII art.",

	Run: func(cmd *cobra.Command, args []string) {
		if !CheckFilePath(cmd, &pFlag) {
			os.Exit(1)
		}

		if !CheckOutputPath(cmd, &oFlag) {
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&pFlag, "path", "p", "", "Path to the file.")
	rootCmd.Flags().StringVarP(&oFlag, "output", "o", ".", "Output folder path. Save the ASCII art to a file.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 

// Returns all the command line arguments.
func GetCommands() (string, string) {
	return pFlag, oFlag
}

// Checks whether the file extension is supported.
func CheckExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".tiff", ".bmp":
		return true
	default:
		return false
	}
}

// Checks whether the file path is valid and the file extension is supported.
func CheckFilePath(cmd *cobra.Command, path *string) bool {
	if *path == "" {
		cmd.Help()
		return false
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		fmt.Printf("The path or file \"%s\" does not exist.\n", *path)
		return false
	}

	if !CheckExtension(*path) {
		fmt.Printf("The file extension is not supported.\n")
		return false
	}

	return true
}

// Checks whether the output path is valid and a directory.
func CheckOutputPath(cmd *cobra.Command, path *string) bool {
	if *path == "" {
        return true
    }

    if *path == "." {
        cwd, err := os.Getwd()
        if err != nil {
            fmt.Printf("Error getting current working directory: %v\n", err)
            return false
        }
        *path = cwd
        return true
    }

    info, err := os.Stat(*path)
    if os.IsNotExist(err) {
        fmt.Printf("The output folder \"%s\" does not exist.\n", *path)
        return false
    } else if err != nil {
        fmt.Printf("Error: %v\n", err)
        return false
    } else if !info.IsDir() {
        fmt.Printf("The output path \"%s\" is not a directory.\n", *path)
        return false
    }

    return true
}