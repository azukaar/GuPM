package ui

import (
	"fmt"
	tm "github.com/buger/goterm"
)

var errorList = make([]string, 0)
var currentLog string
var progress int
var positionToDrawAt int 

func Log(log string) {
	currentLog = log
	draw()
}

func Error(err string) {
	errorList = append(errorList, err)
	draw()
}

func Progress(p int) {
	progress = p
	draw()
}

func moveCursor(x int, y int) {
	fmt.Print("\033[%d;%dH", y, x)
}

func init() {
	positionToDrawAt = 0
}

func draw() {
	moveCursor(1,1)
	fmt.Print("\033[2J")
	fmt.Println(currentLog)
	for _, v := range errorList {
		fmt.Println(v)
	}
}