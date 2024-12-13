package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type Command struct {
	Path   	string
	Output 	string
	Size  	int
}

var cmdFlags Command

var rootCmd = &cobra.Command{
	Use:  "goskii",
	Short: "goskii is a CLI tool to convert images to ASCII art.",

	Run: func(cmd *cobra.Command, args []string) {
		if !checkFilePath(cmd, &cmdFlags.Path) {
			os.Exit(1)
		}

		if !checkOutputPath(cmd, &cmdFlags.Output) {
			os.Exit(1)
		}

		if !checkSize(cmdFlags.Size) {
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&cmdFlags.Path, "path", "p", "", "Path to the file. (Required)")
	rootCmd.Flags().StringVarP(&cmdFlags.Output, "output", "o", "", "Output folder path. Save the ASCII art to a file. ('.' for current directory)")
	rootCmd.Flags().IntVarP(&cmdFlags.Size, "size", "s", 0, "Size of the ASCII art (1 - 100). By default the ASCII art will be scaled to the size of the terminal.")
	rootCmd.MarkPersistentFlagRequired("path")

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 

// Returns all the command line arguments.
func GetCommands() Command {
	return cmdFlags
}

// Checks whether the file extension is supported.
func checkExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".tiff", ".bmp":
		return true
	default:
		return false
	}
}

// Checks whether the file path is valid and the file extension is supported.
func checkFilePath(cmd *cobra.Command, path *string) bool {
	if *path == "" {
		cmd.Help()
		return false
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		fmt.Printf("The path or file \"%s\" does not exist.\n", *path)
		return false
	}

	if !checkExtension(*path) {
		fmt.Printf("The file extension is not supported.\n")
		return false
	}

	return true
}

// Checks whether the output path is valid and a directory.
func checkOutputPath(_ *cobra.Command, path *string) bool {
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

func checkSize(size int) bool {
	if size < 1 || size > 100 {
		fmt.Println("The size should be between 1 and 100.")
		return false
	}

	return true
}