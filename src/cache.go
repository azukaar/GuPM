package main

import (
	"github.com/mitchellh/go-homedir"
	"./ui"
	"./utils"
)


func CacheClear() {
	hdir, errH := homedir.Dir()
	if(errH != nil) {
		ui.Error(errH)
		hdir = "."
	}

	folder := utils.Path(hdir + "/.gupm/cache/")

	utils.RemoveFiles([]string{folder})
}