package main

import (
	"fmt"

	"github.com/JoelVCrasta/goskii/cmd"
	"github.com/JoelVCrasta/goskii/convertor"
)

func main() {
	cmd.Execute()

	pFlag, _ := cmd.GetCommands()
	
	err := convertor.ImageToASCII(pFlag)
	if err != nil {
		fmt.Print(err)
	}
}