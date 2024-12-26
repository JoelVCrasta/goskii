package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	DefaultSize     = 0
    MinSize         = 0
    MaxSize         = 500
    DefaultCharset  = 1
    MinCharset      = 1
    MaxCharset      = 10
	Version 		= "1.0"
)

type Command struct {
	Path   			string
	Output 			string
	Size  			int
	Charset 		int
}
var cmdFlags Command

var rootCmd = &cobra.Command{
	Use:  "goskii",
	Short: "goskii is a CLI tool to convert images to ASCII art.",

	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("showset") {
			showShowset()
			os.Exit(0)
		}

		if cmd.Flags().Changed("version") {
			fmt.Printf("goskii version %s\n", Version)
			os.Exit(0)
		}

		if !checkFilePath(cmd, &cmdFlags.Path) {
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
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&cmdFlags.Path, "path", "p", "","Path to the file. (Required)")
    rootCmd.Flags().StringVarP(&cmdFlags.Output, "output", "o", ".", "Output folder path. Default is current directory.")
    rootCmd.Flags().IntVarP(&cmdFlags.Size, "width", "w", DefaultSize, fmt.Sprintf("Width of the ASCII art (%d - %d). Default adjusts to terminal size.", MinSize, MaxSize))
    rootCmd.Flags().IntVarP(&cmdFlags.Charset, "charset", "c", DefaultCharset, fmt.Sprintf("Character set to use (%d - %d).", MinCharset, MaxCharset))
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
		cmd.PrintErrf("The path to the file is required.\n")
		return false
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		cmd.PrintErrf("The path or file \"%s\" does not exist or not valid.\n", *path)
		return false
	}

	if !checkExtension(*path) {
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

// Checks whether the size is between 0 and 500.
func checkSize(cmd *cobra.Command, size *int) bool {
	if *size < MinSize || *size > MaxSize {
		cmd.PrintErrf("The size should be between %d and %d.", MinSize+1, MaxSize)
		return false
	}

	return true
}

// Checks whether the charset is between 1 and 10.
func checkCharset(cmd *cobra.Command, charset *int) bool {
	if *charset < MinCharset || *charset > MinCharset {
		cmd.PrintErrf("The charset should be between %d and %d.", MinCharset, MaxCharset)
		return false
	}

	return true
}

// Displays all the character sets.
func showShowset() {
	charsets := map[string]string{
        "1": ". : - + * ? # % $ @",
        "2": "@ $ % # ? * + - : .",
        "3": "  ░ ▒ ▓ █",
        "4": "█ ▓ ▒ ░  ",
    }

	charsetsDesc := map[string]string{
		"1": "Lightest to darkest characters. Better for terminal. (Default)",
		"2": "Darkest to lightest characters. Better for light background.",
		"3": "Lightest to darkest block shades. (Unicode)",
		"4": "Darkest to lightest block shades. (Unicode)",
	}

	for key, value := range charsets {
		fmt.Printf("%s\n", charsetsDesc[key])
		fmt.Printf("%s) %s\n\n", key, value)
	}

	fmt.Println("Note: The Unicode characters may not work in all terminals.")
}