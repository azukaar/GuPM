package defaultProvider

import (
	"io/ioutil"
	"encoding/json"
	"../utils"
	"fmt"
)

func GetPackageConfig(entryPoint string) map[string]interface {} {
	var packageConfig map[string]interface{}
	b, err := ioutil.ReadFile(entryPoint)
	if err != nil {
		fmt.Println(err, entryPoint)
	}

	json.Unmarshal([]byte(string(b)), &packageConfig)

	return packageConfig
}

func GetDependencyList(config map[string]interface {}) utils.PackageDepedencyListType {
	return config["dependencies"].(utils.PackageDepedencyListType)
}