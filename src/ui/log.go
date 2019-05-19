package ui

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

var errorList = make([]string, 0)
var debugList = make([]string, 0)
var currentLog string
var currentTitle string
var progress int
var screenWidth int
var positionToDrawAt int 

var redrawNeeded = false

func Title(log string) {
	currentTitle = log
	currentLog = ""
	redrawNeeded = true
}

func Log(log string) {
	currentLog = log
	redrawNeeded = true
}

func Error(err string) {
	errorList = append(errorList, err)
	if(len(errorList) <= 10) {
		Draw()
	}
}

func Debug(err string) {
	debugList = append(debugList, err)
	if(len(debugList) <= 10) {
		Draw()
	}
}

func Progress(p int) {
	progress = p
	redrawNeeded = true
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

	go (func() {
		for _ = range time.Tick(300 * time.Millisecond) {
			if(redrawNeeded) {
				Draw()
			}
		}
	})()
}

func drawTitle() {
	title := color.New(color.FgBlue, color.Bold)
	title.Println("🐶   " + currentTitle)
}

func drawLog() {
	if(currentLog != "") {
		color.Green("✔️   " + currentLog)
	} 
}

func Draw() {
	moveCursor(1,1)
	fmt.Print("\033[2J")
	
	drawTitle()

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
	
	drawLog()

	limit := 0
	for _, v := range errorList {
		if(limit == 10) {
			color.Red("❌❌❌ Too many errors to display...")
			limit++
		} else if(limit < 10) {
			color.Red("❌ " + v)
			limit++
		}
	}

	limit = 0
	for _, v := range debugList {
		if(limit == 10) {
			fmt.Println("Too many debugs...")
			limit++
		} else if(limit < 10) {
			fmt.Println(v)
			limit++
		}
	}
}
