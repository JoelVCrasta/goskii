package cmd

import (
	"fmt"
	"os"

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
		if pFlag == "" {
			cmd.Help()
			os.Exit(1)
		}

		// Check if the path exists.
		if _, err := os.Stat(pFlag); os.IsNotExist(err) {
			fmt.Printf("The path or file \"%s\" does not exist.\n", pFlag)
			os.Exit(1)
		}

		// Check if the output folder exists.
		if oFlag != "" {
			// No-op if -o not provided. 
		} else if oFlag == "." {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			oFlag = cwd
		} else {
			if info, err := os.Stat(oFlag); os.IsNotExist(err) {
				fmt.Printf("The output folder \"%s\" does not exist.\n", oFlag)
				os.Exit(1)
			} else if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			} else if !info.IsDir() {
				fmt.Printf("The output path \"%s\" is not a directory.\n", oFlag)
				os.Exit(1)
			}
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

func GetPaths() (string, string) {
	return pFlag, oFlag
}