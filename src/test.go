package main

import (
	"regexp"
	"os"
	"strconv"
	"./utils"
	"./ui"
	"./jsVm"
)

func RunTest(path string) {
	_ = jsVm.Run
	files := utils.RecursiveFileWalkDir(path)
	i := 0
	for _, file := range files {
		isTest, _ := regexp.MatchString(`\.test\.gs$`, file)
		if(isTest) {
			os.MkdirAll(".tmp_test_gupm", os.ModePerm)
			utils.CopyFiles([]string{file}, ".tmp_test_gupm/"+strconv.Itoa(i)+".gs")
			errch := os.Chdir(".tmp_test_gupm")
			if(errch != nil) {
				ui.Error("Could'n execute", file)
				Exit(1)
			}
			_, err := jsVm.Run(strconv.Itoa(i)+".gs", make(map[string]interface {}))
			os.Chdir("..")
			utils.RemoveFiles([]string{".tmp_test_gupm"})
			if(err != nil) {
				ui.Error("Test execution failed:", file)
				Exit(1)
			}
			i++
		}
	}

	ui.Title("Test passed! ❤️")
}