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
	
	vm.Set("getDependency", func(call otto.FunctionCall) otto.Value {
		providerName, _ := call.Argument(0).ToString()
		name, _ := call.Argument(1).ToString()
		version, _ := call.Argument(2).ToString()
		url, _ := call.Argument(3).ToString()
		path := "./cache/" + providerName + "/" + name + "/" + version
		
		res, errorGD := GetDependency(providerName, name, version, url, path)
		if(errorGD != nil) {
			fmt.Println(errorGD)
		}
		postrest, errorPGD := PostGetDependency(providerName, name, version, url, path, res)
		if(errorPGD != nil) {
			fmt.Println(errorPGD)
		}

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

func ExpandDependency(dependency map[string]interface {}) (map[string]interface {}, error) {
	var file, _ = utils.FileExists(ProviderPath + "/ExpandDependency.js")
	if(file) {
		input := make(map[string]interface {})
		input["Dependency"] = dependency

		res, err :=  run(ProviderPath + "/ExpandDependency.js", input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		deps, _ := resObj.(map[string]interface{})["dependencies"].(otto.Value).Export()
		resObj.(map[string]interface{})["dependencies"] = deps
		return resObj.(map[string]interface {}), err1
	} else {
		return nil, nil
	}
}

func GetDependency(provider string, name string, version string, url string, path string) (string, error) {
	var file, _ = utils.FileExists(ProviderPath + "/GetDependency.js")
	if(file) {
		input := make(map[string]interface {})
		input["Provider"] = provider
		input["Name"] = name
		input["Version"] = version
		input["Url"] = url
		input["Path"] = path

		res, err :=  run(ProviderPath + "/GetDependency.js", input)
		if(err != nil) {
			return "", err
		}

		resStr, err1 := res.ToString()
		return resStr, err1
	} else {
		return defaultProvider.GetDependency(provider, name, version, url, path)
	}
}

func PostGetDependency(provider string, name string, version string, url string, path string, result string) (string, error) {
	var file, _ = utils.FileExists(ProviderPath + "/PostGetDependency.js")
	if(file) {
		input := make(map[string]interface {})
		input["Provider"] = provider
		input["Name"] = name
		input["Version"] = version
		input["Url"] = url
		input["Path"] = path
		input["Result"] = result

		res, err :=  run(ProviderPath + "/PostGetDependency.js", input)
		if(err != nil) {
			return "", err
		}

		resStr, err1 := res.ToString()
		return resStr, err1
	} else {
		return defaultProvider.PostGetDependency(provider, name, version, url, path, result)
	}
}
