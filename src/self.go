package main

import (
	"./utils"
	"./ui"
	"runtime"
	"os/exec"
	"io/ioutil"
	"os"
)

// TODO: implemement script-free upgrade for all OSes

func SelfUpgrade() {
	SelfUninstall()

	if runtime.GOOS != "windows" {
		script := utils.HttpGet("https://azukaar.github.io/GuPM/install.sh")
		ioutil.WriteFile("TEMP_install.sh", []byte(script), os.ModePerm)
		cmd := exec.Command("/bin/sh", "TEMP_install.sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
		utils.RemoveFiles([]string{"TEMP_install.sh"})
	} else {
		ui.Error("Upgrade not fully implememnted on windows yet. Please execute windows_installer.exe again")
	}
}

func SelfUninstall() {
	utils.RemoveFiles([]string{utils.DIRNAME()})
	
	if runtime.GOOS != "windows" {
		utils.RemoveFiles([]string{"/usr/local/bin/g", "/bin/g"})
	}
}