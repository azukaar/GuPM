package main

import (
	"regexp"
	"time"
	"./utils"
	"./jsVm"
)

func RunTest(path string) error {
	files := utils.RecursiveFileWalkDir(path)

	for _, file := range files {
		isTest, _ := regexp.MatchString(`\.test\.gs$`, file)
		if(isTest) {
			jsVm.Run(file, make(map[string]interface {}))
		}
	}

	return nil
}