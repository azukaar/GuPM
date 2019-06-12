package main

import (
	"./provider"
	"./ui"
	// "fmt"
)

func Bootstrap(path string) error {	
	err := provider.InitProvider(Provider)
	if(err != nil) {
		return err
	}
	ui.Title("Bootstrap project")
	return provider.Bootstrap(path)
}