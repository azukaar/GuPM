package main

import (
	"./provider"
	"./ui"
	// "fmt"
)

func Bootstrap(path string) error {
	err := provider.InitProvider(Provider)
	if err != nil {
		return err
	}
	ui.Title("Bootstrap project")
	errBoot := provider.Bootstrap(path)
	if errBoot != nil {
		return errBoot
	} else {
		ui.Title("Bootstrap done! ❤️")
	}
	return nil
}
