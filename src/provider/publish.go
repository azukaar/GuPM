package provider

import (
	"../defaultProvider"
	"../jsVm"
	"errors"
	"../utils"
	// "fmt"
)

func Publish(path string) error {
	if(Provider != "gupm") {
		var file = utils.FileExists(utils.Path(ProviderPath + "/publish.gs"))
		if(file) {
			input := make(map[string]interface {})
			input["Path"] = path
			_, err :=  jsVm.Run(utils.Path(ProviderPath + "/publish.gs"), input)
			return err
		} else {
			return errors.New("Provider doesn't have any publish function. Please use 'g publish' to use the default publish.")
		}
	} else {
		return defaultProvider.Publish(path)
	}
}