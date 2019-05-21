package ui

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"time"
	"sync"
	// "io"
)

var errorList = make([]string, 0)
var debugList = make([]string, 0)
var currentLog string
var currentTitle string
var progress int
var screenWidth int
var positionToDrawAt int 
var logBox = uilive.New()

var lock = sync.RWMutex{}

var redrawNeeded = false
var running = true

func Title(log string) {
    _ = color.Green
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
		redrawNeeded = true
		// Draw()
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

    logBox.Start()
	
	go (func() {
		for _ = range time.Tick(200 * time.Millisecond) {
			if(running){
				Draw()
			}
		}
	})()
}

func drawTitle() string {
	if(currentTitle != "") {
		title := color.New(color.FgBlue, color.Bold)
		return title.Sprintln("🐶   " + currentTitle)
	} else {
		return ""
	}
}

func drawLog() string {
	if(currentLog != "") {
		log := color.New(color.FgGreen)
		return log.Sprintln("✓   " + currentLog)
	} else {
		return ""
	}
}

func Stop() {
	redrawNeeded = true
	running = false
	Draw()
}

func Draw() {
	if(!redrawNeeded) {
		return;
	}

	result := ""

	result += drawTitle()

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
	
	result += drawLog()

	errorColor := color.New(color.FgRed)
	limit := 0
	for _, v := range errorList {
		_ = v
		if(limit == 10) {
			result += errorColor.Sprintln("❌❌❌ Too many errors to display...")
			limit++
			} else if(limit < 10) {
			result += errorColor.Sprintln("❌ " + v)
			limit++
		}
	}

	limit = 0
	for _, v := range debugList {
		_ = v
		if(limit == 10) {
			result += "Too many debugs..."
			limit++
		} else if(limit < 10) {
			result += v
			limit++
		}
	}

	lock.Lock()
	if(running) {
		fmt.Fprintf(logBox, result)
	} else {
		fmt.Fprintf(logBox, "\n")
		logBox.Stop() 
		fmt.Println(result)
	}
	redrawNeeded = false
	lock.Unlock()
}
