package main

import (
	"./provider"
	"./ui"
)

func Publish(path string) error {
	err := provider.InitProvider(Provider)
	if err != nil {
		return err
	}
	ui.Title("Publishing package...")
	errPub := provider.Publish(path)
	if errPub != nil {
		return errPub
	} else {
		ui.Title("Publish done! ❤️")
	}
	return nil
}
