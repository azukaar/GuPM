package main

import (
	"regexp"
	"time"
	"./utils"
	"./jsVm"
)

func RunTest(path string) error {
	files := utils.RecursiveFileWalkDir(path)
	nbRunning := 0

	for _, file := range files {
		for nbRunning > 15 {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}

		isTest, _ := regexp.MatchString(`\.test\.gs$`, file)
		if(isTest) {
			jsVm.Run(file, make(map[string]interface {}))
		}
	}

	return nil
}