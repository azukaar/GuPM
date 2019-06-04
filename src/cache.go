package main

import (
	"github.com/mitchellh/go-homedir"
	"./ui"
	"./utils"
)


func CacheClear() {
	hdir, errH := homedir.Dir()
	if(errH != nil) {
		ui.Error(errH.Error())
		hdir = "."
	}

	folder := hdir + "/.gupm/cache/"

	utils.RemoveFiles([]string{folder})
}