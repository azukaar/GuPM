package ui

import (
	"fmt"
	"github.com/fatih/color"
)

var errorList = make([]string, 0)
var debugList = make([]string, 0)
var currentLog string
var currentTitle string
var progress int
var screenWidth int
var positionToDrawAt int 

func Title(log string) {
	currentTitle = log
	currentLog = ""
	draw()
}

func Log(log string) {
	currentLog = log
	draw()
}

func Error(err string) {
	errorList = append(errorList, err)
	draw()
}

func Debug(err string) {
	debugList = append(debugList, err)
	draw()
}

func Progress(p int) {
	progress = p
	draw()
}

// https://github.com/ahmetb/go-cursor/blob/master/cursor.go

var Esc = "\x1b"

func escape(format string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", Esc, fmt.Sprintf(format, args...))
}

func moveCursor(x int, y int) {
	escape("[%d;%dH", x, y)
}

func init() {
	positionToDrawAt = 0
}

func draw() {
	moveCursor(1,1)
	fmt.Print("\033[2J")
	
	title := color.New(color.FgBlue, color.Bold)
	title.Println("🐶   " + currentTitle)
	
	if(progress > 0) {
		fmt.Print("📦📦")
		for i := 0; i < 20; i++  {
			if(i == progress / 5) {
				fmt.Print("🐕")
			} else {
				fmt.Print("-")
			}
		}
		fmt.Println("🏠")
	}

	if(currentLog != "") {
		color.Green("✔️   " + currentLog)
	} 

	for _, v := range errorList {
		color.Red("❌ " + v)
	}

	for _, v := range debugList {
		fmt.Println(v)
	}
}