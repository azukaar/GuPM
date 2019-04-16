package provider

import (
	"fmt"
	"../utils"
	"os"
	"net/http"
	"io/ioutil"
	"github.com/robertkrimen/otto"
	"../defaultProvider"
	"github.com/Masterminds/semver"
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

func httpGet(url string) []byte {
	resp, httperr := http.Get(url)
	if httperr != nil {
		fmt.Println("Error trying to dl file ", url, " trying again. Check your network.")
		return httpGet(url)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	
	return body
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

	vm.Set("httpGet", func(call otto.FunctionCall) otto.Value {
		url, _ := call.Argument(0).ToString()
		res := httpGet(url)
		result, _ :=  vm.ToValue(utils.StringToJSON(string(res)))
		return result
	})

	vm.Set("semverInRange", func(call otto.FunctionCall) otto.Value {
		rangeStr, _ := call.Argument(0).ToString()
		version, _ := call.Argument(1).ToString()
		rangeVer, _ := semver.NewConstraint(rangeStr)
		sver, _ := semver.NewVersion(version)
		value := rangeVer.Check(sver)
		result, _ :=  vm.ToValue(value)
		return result
	})

	vm.Set("semverLatestInRange", func(call otto.FunctionCall) otto.Value {
		rangeStr, _ := call.Argument(0).ToString()
		versionList, _ := call.Argument(1).Export()
		var version string
		rangeVer, _ := semver.NewConstraint(rangeStr)

		for _, verCand := range versionList.([]string) {
			sver, err := semver.NewVersion(verCand)
			if err != nil {
				fmt.Println("!", err)
			}
	
			if(rangeVer.Check(sver)) {
				version = verCand
			}
		}
		if(version != "") {
			result, _ :=  vm.ToValue(version)
			return result
		} else {
			return otto.UndefinedValue()
		}
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

func ExpandDependency(depedency map[string]interface {}) (map[string]interface {}, error) {
	var file, _ = fileExists(ProviderPath + "/ExpandDependency.js")
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
