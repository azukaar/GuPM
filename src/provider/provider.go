package provider

import (
	"fmt"
	"../utils"
	"os"
	"io/ioutil"
	"github.com/robertkrimen/otto"
	"../defaultProvider"
)

var Provider string
var ProviderPath string
var ProviderConfig *gupmEntryPoint
var scriptCache = make(map[string]string)

func fileExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func run(path string, input map[string]interface {}) (otto.Value, error) {
	var err error
	var ret otto.Value

	if(scriptCache[path] == "") {
		file, err := ioutil.ReadFile(path)
		if(err != nil) {
			return otto.UndefinedValue(),  err
		}
		scriptCache[path] = string(file)
	}

	vm := otto.New()
	for varName, varValue := range input {
		vm.Set(varName, varValue)
	}

	ret, err = vm.Run(scriptCache[path])

	if(err != nil) {
		return otto.UndefinedValue(),  err
	}

	return ret, nil
}

func InitProvider(provider string) error {
	Provider = provider
	ProviderPath = "plugins/provider-" + Provider
	fmt.Println("Reading provider config for", Provider)
	ProviderConfig = new(gupmEntryPoint) 
	err := utils.ReadJSON(ProviderPath + "/gupm.json", ProviderConfig)
	if(err != nil) {
		return err
	} else {
		fmt.Println("Initialisation OK for", ProviderConfig.Name);
	}

	return nil
}

func GetEntryPoint() string {
	fmt.Println(ProviderConfig.Name)
	return ProviderConfig.Config.Default.Entrypoint
}

func GetPackageConfig() (utils.Json, error) {
	var file, _ = fileExists(ProviderPath + "/GetPackageConfig.js")
	if(file) {
		input := make(map[string]interface {})		
		res, err :=  run(ProviderPath + "/GetPackageConfig.js", input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		return resObj.(utils.Json), err1
	} else {
		return defaultProvider.GetPackageConfig(GetEntryPoint()), nil
	}
}

func PostGetPackageConfig(config utils.Json) (utils.Json, error) {
	var file, _ = fileExists(ProviderPath + "/PostGetPackageConfig.js")
	if(file) {
		input := make(map[string]interface {})
		input["PackageConfig"] = config
		
		res, err :=  run(ProviderPath + "/PostGetPackageConfig.js", input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		return resObj.(utils.Json), err1
	} else {
		return config, nil
	}
}

func GetDependencyList(config utils.Json) ([]map[string]interface {}, error) {
	var file, _ = fileExists(ProviderPath + "/GetDependencyList.js")
	if(file) {
		input := make(map[string]interface {})
		input["PackageConfig"] = config
		
		res, err :=  run(ProviderPath + "/GetDependencyList.js", input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		return resObj.([]map[string]interface {}), err1
	} else {
		return defaultProvider.GetDependencyList(config), nil
	}
}

func ExpandDependency(depedency map[string]interface {}) ([]interface {}, error) {
	var file, _ = fileExists(ProviderPath + "/ExpandDependency.js")
	if(file) {
		input := make(map[string]interface {})
		input["Depedency"] = depedency

		res, err :=  run(ProviderPath + "/ExpandDependency.js", input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		return resObj.([]interface {}), err1
	} else {
		return nil, nil
	}
}
