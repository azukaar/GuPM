package provider

import (
	"../defaultProvider"
	"../utils"
	"../jsVm"
	"../ui"
	"sync"
)

var Provider string
var ProviderPath string
var providerConfigCache = make(map[string]*GupmEntryPoint)
var linkHasErrored = false
var pConfigLock = sync.RWMutex{}

func InitProvider(provider string) error {
	Provider = provider
	
	if(provider == "gupm") {
		ProviderPath = utils.DIRNAME()
	} else {
		ProviderPath = utils.DIRNAME() + utils.Path("/plugins/provider-" + Provider)
	}

	if(Provider != "") {
		providerConfig := GetProviderConfig(Provider) 
		ui.Log("Initialisation OK for " + providerConfig.Name);
	} else {
		providerConfig := GetProviderConfig("gupm") 
		ui.Log("Initialisation OK for " + providerConfig.Name);
	}
	return nil
}

func GetProviderConfig(providerName string) *GupmEntryPoint {
	var providerConfigPath string

	if(providerName == "gupm") {
		providerConfigPath = utils.DIRNAME() + utils.Path("/gupm.json")
	} else {
		providerConfigPath = utils.DIRNAME() + utils.Path("/plugins/provider-" + providerName + "/gupm.json")
	}

	pConfigLock.Lock()
	if(providerConfigCache[providerName] == nil) {
		config := new(GupmEntryPoint)
		err := utils.ReadJSON(providerConfigPath, config)
		if(err != nil) {
			ui.Error(err.Error())
			return nil
		}

		providerConfigCache[providerName] = config
		
		pConfigLock.Unlock()
		return config
	} else {
		config := providerConfigCache[providerName]
		pConfigLock.Unlock()
		return config
	}
}

func GetPackageConfig() (utils.Json, error) {
	var file = utils.FileExists(ProviderPath + utils.Path("/getPackageConfig.gs"))
	if(file) {
		input := make(map[string]interface {})		
		res, err :=  jsVm.Run(ProviderPath + utils.Path("/getPackageConfig.gs"), input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		return resObj.(utils.Json), err1
	} else {
		return defaultProvider.GetPackageConfig(GetProviderConfig(Provider).Config.Default.Entrypoint), nil
	}
}

func PostGetPackageConfig(config utils.Json) (utils.Json, error) {
	var file = utils.FileExists(ProviderPath + utils.Path("/postGetPackageConfig.gs"))
	if(file) {
		input := make(map[string]interface {})
		input["PackageConfig"] = config
		
		res, err :=  jsVm.Run(ProviderPath + utils.Path("/postGetPackageConfig.gs"), input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		return resObj.(utils.Json), err1
	} else {
		return config, nil
	}
}

func SaveDependencyList(depList []map[string]interface {}) error {
	var file = utils.FileExists(ProviderPath + utils.Path("/saveDependencyList.gs"))
	if(file) {
		input := make(map[string]interface {})
		input["Dependencies"] = depList
		
		_, err :=  jsVm.Run(ProviderPath + utils.Path("/saveDependencyList.gs"), input)
		if(err != nil) {
			return err
		}

		return nil
	} else {
		return defaultProvider.SaveDependencyList(depList)
	}
}

func GetDependencyList(config utils.Json) ([]map[string]interface {}, error) {
	var file = utils.FileExists(ProviderPath + utils.Path("/getDependencyList.gs"))
	if(file) {
		input := make(map[string]interface {})
		input["PackageConfig"] = config
		
		res, err :=  jsVm.Run(ProviderPath + utils.Path("/getDependencyList.gs"), input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		resMap, ok := resObj.([]map[string]interface {})

		if(ok) {
			return resMap, err1
		} else {
			return make([]map[string]interface {}, 0), err1
		}
	} else {
		return defaultProvider.GetDependencyList(config), nil
	}
}

func ResolveDependencyLocation(dependency map[string]interface {}) (map[string]interface {}, error) {
	depProviderPath := utils.DIRNAME() + utils.Path("/plugins/provider-" + dependency["provider"].(string))
	var file = utils.FileExists(depProviderPath + utils.Path("/resolveDependencyLocation.gs"))
	if(dependency["provider"].(string) != "gupm" && file) {
		input := make(map[string]interface {})
		input["Dependency"] = dependency

		res, err :=  jsVm.Run(depProviderPath + utils.Path("/resolveDependencyLocation.gs"), input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		
		if(resObj == nil) {
			ui.Error("ERROR Failed to resolve" + dependency["name"].(string) + "Trying again.")
			return ResolveDependencyLocation(dependency)
		}
		return resObj.(map[string]interface {}), err1
	} else {
		return nil, nil
	}
}

func ExpandDependency(dependency map[string]interface {}) (map[string]interface {}, error) {
	depProviderPath := utils.DIRNAME() + utils.Path("/plugins/provider-" + dependency["provider"].(string))
	var file = utils.FileExists(depProviderPath + utils.Path("/expandDependency.gs"))
	if(dependency["provider"].(string) != "gupm" && file) {
		input := make(map[string]interface {})
		input["Dependency"] = dependency

		res, err :=  jsVm.Run(depProviderPath + utils.Path("/expandDependency.gs"), input)
		if(err != nil) {
			return nil, err
		}

		toExport, _ := res.Export()
		resObj := jsVm.JsonExport(toExport).(map[string] interface {})

		if(resObj == nil) {
			ui.Error("ERROR Failed to resolve" + dependency["name"].(string) + ". Trying again.")
			return ExpandDependency(dependency)
		}

		return resObj, nil
	} else {
		return defaultProvider.ExpandDependency(dependency)
	}
}

func GetDependency(provider string, name string, version string, url string, path string) (string, error) {
	depProviderPath := utils.DIRNAME() + utils.Path("/plugins/provider-" + provider)
	var file = utils.FileExists(depProviderPath + utils.Path("/getDependency.gs"))
	if(provider != "gupm" && file) {
		input := make(map[string]interface {})
		input["Provider"] = provider
		input["Name"] = name
		input["Version"] = version
		input["Url"] = url
		input["Path"] = path

		res, err :=  jsVm.Run(depProviderPath + utils.Path("/getDependency.gs"), input)
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
	depProviderPath := utils.DIRNAME() + utils.Path("/plugins/provider-" + provider)
	var file = utils.FileExists(depProviderPath + utils.Path("/postGetDependency.gs"))
	if(provider != "gupm" && file) {
		input := make(map[string]interface {})
		input["Provider"] = provider
		input["Name"] = name
		input["Version"] = version
		input["Url"] = url
		input["Path"] = path
		input["Result"] = result

		res, err :=  jsVm.Run(depProviderPath + utils.Path("/postGetDependency.gs"), input)
		if(err != nil) {
			return "", err
		}

		resStr, err1 := res.ToString()
		return resStr, err1
	} else {
		return defaultProvider.PostGetDependency(provider, name, version, url, path, result)
	}
}

