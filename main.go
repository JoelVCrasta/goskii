package main

import (
	"fmt"

	"github.com/JoelVCrasta/goskii/cmd"
	"github.com/JoelVCrasta/goskii/convertor"
	"github.com/JoelVCrasta/goskii/utils"
)

func main() {
	cmd.Execute()
	cmdFlags := cmd.GetCommands()
	ftype := cmd.GetFileType()

	if cmdFlags.Path != "" {
		if ftype == 0 {
			err := convertor.ImageToASCII(cmdFlags)
			if err != nil {
				fmt.Print(err)
			}
		} else if ftype == 1 || ftype == 3 {
			err := convertor.VideoToASCII(cmdFlags)
			if err != nil {
				fmt.Print(err)
			}
		} else {
			fmt.Print("Invalid file type")
		}
	} else if cmdFlags.Render != "" {
		utils.Render(cmdFlags.Render, cmdFlags.Fps)
	}
}