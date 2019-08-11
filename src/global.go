package main

import (
	"os"
	"path/filepath"

	"./ui"
	"./utils"
)

func installBins(path string) {
	var GLOBAL = utils.Path(utils.HOMEDIR(".") + "/.gupm/global")
	bins := utils.RecursiveFileWalkDir(GLOBAL + "/.bin")

	for _, bin := range bins {
		name := filepath.Base(bin)
		if !utils.FileExists(path + "/" + name) {
			os.Symlink(bin, path+"/"+name)
		}
	}
}

func GlobalAdd(rls []string) {
	ui.Title("Installing global dependency...")
	var GLOBAL = utils.Path(utils.HOMEDIR(".") + "/.gupm/global")

	if !utils.FileExists(GLOBAL + utils.Path("/gupm.json")) {
		os.MkdirAll(GLOBAL, os.ModePerm)
		utils.WriteFile(GLOBAL+utils.Path("/gupm.json"), "{}")
	}

	ui.Log("Installing...")
	AddDependency(GLOBAL, rls)
	InstallProject(GLOBAL)

	ui.Log("Add binaries...")
	if utils.OSNAME() != "windows" {
		if utils.FileExists("/usr/local/bin/") {
			installBins("/usr/local/bin/")
		} else {
			installBins("/usr/bin/")
		}
	} else {
		ui.Error("Global Installation not supported on Windows yet. Please add .gupm/global/.bin to your PATH")
	}
}

func GlobalDelete(rls []string) {
	var GLOBAL = utils.Path(utils.HOMEDIR(".") + "/.gupm/global")

	if !utils.FileExists(GLOBAL + utils.Path("/gupm.json")) {
		os.MkdirAll(GLOBAL, os.ModePerm)
		utils.WriteFile(GLOBAL+utils.Path("/gupm.json"), "{}")
	}

	RemoveDependency(GLOBAL, rls)
}
