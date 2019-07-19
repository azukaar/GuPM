package main

import (
	"./provider"
	"./ui"
)

func Publish(path string, namespace string) error {
	err := provider.InitProvider(Provider)
	if err != nil {
		return err
	}
	ui.Title("Publishing package...")
	errPub := provider.Publish(path, namespace)
	if errPub != nil {
		return errPub
	} else {
		ui.Title("Publish done! ❤️")
	}
	return nil
}
