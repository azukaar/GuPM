package provider

import (
	"../defaultProvider"
	"../jsVm"
	"../utils"
	"errors"
	// "fmt"
)

func Bootstrap(path string) error {
	if Provider != "gupm" {
		var file = utils.FileExists(utils.Path(ProviderPath + "/bootstrap.gs"))
		if file {
			input := make(map[string]interface{})
			input["Path"] = path
			_, err := jsVm.Run(utils.Path(ProviderPath+"/bootstrap.gs"), input)
			if err != nil {
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
