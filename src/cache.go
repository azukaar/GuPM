package main

import (
	"./ui"
	"./utils"
	"github.com/mitchellh/go-homedir"
)

func CacheClear() {
	hdir, errH := homedir.Dir()
	if errH != nil {
		ui.Error(errH)
		hdir = "."
	}

	folder := utils.Path(hdir + "/.gupm/cache/")

	utils.RemoveFiles([]string{folder})
}
