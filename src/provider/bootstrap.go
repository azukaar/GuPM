package provider

import (
	"../defaultProvider"
	"errors"
	"../jsVm"
	"../utils"
	// "fmt"
)

func Bootstrap(path string) error {
	if(Provider != "gupm") {
		var file = utils.FileExists(ProviderPath + "/Bootstrap.js")
		if(file) {
			input := make(map[string]interface {})		
			_, err :=  jsVm.Run(ProviderPath + "/Bootstrap.js", input)
			if(err != nil) {
				return err
			}
		} else {
			return errors.New("Provider doesn't have any bootstrap function. Please use 'g bootstrap' to use the default bootstrap.")
		}
	} else {
		defaultProvider.Bootstrap(path)
	}

	return nil
}