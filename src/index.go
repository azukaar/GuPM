package main

import (
	"os"
	"./utils"
	"./ui"
	"./jsVm"
	"strings"
	"path/filepath"
	"regexp"
	"fmt"
	"runtime"
)

var Provider string

func ScriptExists(path string) string {
	if (utils.FileExists(path + ".gs")) {
		return path + ".gs"
	} else if(utils.FileExists(path) && !utils.IsDirectory(path)) {
		return path
	} else {
		return ""
	}
}

func executeFile(path string, args Arguments) {
	_, err := jsVm.Run(path, args.AsJson())
	if(err != nil) {
	  ui.Error("File execution failed")
		ui.Error(err)
		os.Exit(1)
	}
}

func binFile(name string, args []string) {
	path := utils.Path("./.bin/"+name)
	realPath, _ := filepath.EvalSymlinks(path)
	utils.ExecCommand(realPath, args)
}

func main() {
	if runtime.GOOS == "darwin" {
		utils.ExecCommand("ulimit", []string{"-n", "2048"})
	}
	
	binFolder := make(map[string]bool)

	if(utils.FileExists(".bin")) {
		for _, file := range utils.ReadDir(".bin") {
			binFolder[file.Name()] = true
		}
	}
	
	c, args := GetArgs(os.Args[1:])

	if(utils.FileExists(".gupm_rc.gs")) {
		executeFile(".gupm_rc.gs", args)
	}
	
	aliases := map[string] interface{}{}
	if(utils.FileExists("gupm.json")) {
		packageConfig, errConfig := utils.ReadGupmJson("gupm.json")
		if(errConfig != nil) {
			ui.Error(errConfig)
		} else {
			aliases = packageConfig.Cli.Aliases
		}
	}

	script := ScriptExists(c)
	if didExec, err := ExecCli(c, args); didExec {
		if(err != nil) {
			ui.Error(err);
		}
		if (script != "") {
			executeFile(script, args)
		}
	} else if (c == "env" || c == "e") {
		toProcess := os.Args[2:]
		re := regexp.MustCompile(`([\w\-\_]+)=([\w\-\_]+)`)
		isEnv := re.FindAllStringSubmatch(toProcess[0], -1)
		for(isEnv != nil) {
			name := isEnv[0][1]
			value := isEnv[0][2]
			os.Setenv(name, value)
			toProcess = toProcess[1:]
			isEnv = re.FindAllStringSubmatch(toProcess[0], -1)
		}
		utils.ExecCommand(toProcess[0], toProcess[1:])
	} else if (aliases[c] != nil) {
		commands := strings.Split(aliases[c].(string), ";")
		for _, command := range commands {
			commandList := strings.Split(command, " ")
			utils.ExecCommand(commandList[0], append(commandList[1:], os.Args[2:]...))
		}
	} else if (binFolder[c] == true) {
		binFile(c, os.Args[2:])	
	} else if (script != "") {
		executeFile(script, args)
	} else if (c == "") {
		fmt.Println("Welcome to GuPM version 1.0.0 \ntry 'g help' for a list of commands. Try 'g filename' to execute a file.")
	} else {
		fmt.Println("Command not found. Try 'g help' or check filename.")
	}

	ui.Stop()
}