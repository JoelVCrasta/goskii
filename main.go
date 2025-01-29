package main

import (
	"fmt"
	"os"

	"github.com/JoelVCrasta/goskii/cmd"
	"github.com/JoelVCrasta/goskii/convertor"
	"github.com/JoelVCrasta/goskii/utils"
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

	ftype := cmd.GetFileType()

	if (cmdFlags.Path != "") {
		if ftype == 0 {
			err := convertor.ImageToASCII(cmdFlags)
			if err != nil {
				fmt.Print(err)
			}
		} else if ftype == 1 {
			err := convertor.VideoToASCII(cmdFlags)
			if err != nil {
				fmt.Print(err)
			}
		}
	} else if (cmdFlags.Render != "") {
		utils.Render(cmdFlags.Render)
	}
}