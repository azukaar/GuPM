package main

import (
	"github.com/robertkrimen/otto"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	vm := otto.New()
	vm.Run(`
		abc = 2 + 2;
		console.log("The value of abc is " + abc); // 4
	`)
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}