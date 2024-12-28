package main

import (
	"fmt"
	"os"

	"github.com/JoelVCrasta/goskii/cmd"
	"github.com/JoelVCrasta/goskii/convertor"
	"github.com/JoelVCrasta/goskii/imageutils"
)

func main() {
	cmd.Execute()
	cmdFlags := cmd.GetCommands()

	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-h" || arg == "--help" {
				return
			}
		}
	}

	if (cmdFlags.Path != "") {
		err := convertor.ImageToASCII(cmdFlags)
		if err != nil {
			fmt.Print(err)
		}
	} else if (cmdFlags.Render != "") {
		imageutils.Render(cmdFlags.Render)
	}
	
}