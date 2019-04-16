package provider

import (
	"fmt"
	"../defaultProvider"
	"../utils"
	"../jsVm"
	"io/ioutil"
	"github.com/robertkrimen/otto"
)

var Provider string
var ProviderPath string
var ProviderConfig *gupmEntryPoint
var scriptCache = make(map[string]string)

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
	jsVm.Setup(vm)
	
	vm.Set("getDepedency", func(call otto.FunctionCall) otto.Value {
		providerName, _ := call.Argument(0).ToString()
		name, _ := call.Argument(1).ToString()
		version, _ := call.Argument(2).ToString()
		url, _ := call.Argument(3).ToString()
		path := "./cache/" + providerName + "/" + name + "/" + version
		
		res, _ := GetDepedency(providerName, name, version, url, path)
		postrest, _ := PostGetDepedency(providerName, name, version, url, path, res)

		result, _ :=  vm.ToValue(postrest)
		return result
	})

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
	var file, _ = utils.FileExists(ProviderPath + "/GetPackageConfig.js")
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
	var file, _ = utils.FileExists(ProviderPath + "/PostGetPackageConfig.js")
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
	var file, _ = utils.FileExists(ProviderPath + "/GetDependencyList.js")
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

func ExpandDependency(depedency map[string]interface {}) (map[string]interface {}, error) {
	var file, _ = utils.FileExists(ProviderPath + "/ExpandDependency.js")
	if(file) {
		input := make(map[string]interface {})
		input["Depedency"] = depedency

		res, err :=  run(ProviderPath + "/ExpandDependency.js", input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		return resObj.(map[string]interface {}), err1
	} else {
		return nil, nil
	}
}

func GetDepedency(provider string, name string, version string, url string, path string) (string, error) {
	var file, _ = utils.FileExists(ProviderPath + "/GetDepedency.js")
	if(file) {
		input := make(map[string]interface {})
		input["Provider"] = provider
		input["Name"] = name
		input["Version"] = version
		input["Url"] = url
		input["Path"] = path

		res, err :=  run(ProviderPath + "/GetDepedency.js", input)
		if(err != nil) {
			return "", err
		}

		resStr, err1 := res.ToString()
		return resStr, err1
	} else {
		return defaultProvider.GetDepedency(provider, name, version, url, path)
	}
}

func PostGetDepedency(provider string, name string, version string, url string, path string, result string) (string, error) {
	var file, _ = utils.FileExists(ProviderPath + "/PostGetDepedency.js")
	if(file) {
		input := make(map[string]interface {})
		input["Provider"] = provider
		input["Name"] = name
		input["Version"] = version
		input["Url"] = url
		input["Path"] = path

		res, err :=  run(ProviderPath + "/PostGetDepedency.js", input)
		if(err != nil) {
			return "", err
		}

		resStr, err1 := res.ToString()
		return resStr, err1
	} else {
		return defaultProvider.PostGetDepedency(provider, name, version, url, path, result)
	}
}
