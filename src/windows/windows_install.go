package main

import (
	"../ui"
	"../utils"
	"fmt"
)

func main() {
	arch := utils.HttpGet("https://azukaar.github.io/GuPM/gupm_windows.tar.gz")
	files, _ := utils.Untar(string(arch))
	path := ui.WaitForInput("Where do you want to save GuPM? (default C:\\)")
	if path == "" {
		path = "C:\\"
	}
	files.SaveAt(path)
	fmt.Println("GuPM saved in " + path + ". dont forget to add it to your PATH")
}
