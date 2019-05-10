package provider

import (
	"fmt"
	"../defaultProvider"
	"../utils"
	"../jsVm"
	"os"
)

var Provider string
var ProviderPath string
var providerConfigCache = make(map[string]*GupmEntryPoint)

func InitProvider(provider string) error {
	Provider = provider
	
	if(provider == "gupm") {
		ProviderPath = utils.DIRNAME()
	} else {
		ProviderPath = utils.DIRNAME() + "/plugins/provider-" + Provider
	}

	if(Provider != "") {
		providerConfig := GetProviderConfig(Provider) 
		fmt.Println("Initialisation OK for", providerConfig.Name);
	} else {
		providerConfig := GetProviderConfig("gupm") 
		fmt.Println("Initialisation OK for", providerConfig.Name);
	}
	return nil
}

func GetProviderConfig(providerName string) *GupmEntryPoint {
	var providerConfigPath string

	if(providerName == "gupm") {
		providerConfigPath = utils.DIRNAME() + "/gupm.json"
	} else {
		providerConfigPath = utils.DIRNAME() + "/plugins/provider-" + providerName + "/gupm.json"
	}

	if(providerConfigCache[providerName] == nil) {
		config := new(GupmEntryPoint)
		err := utils.ReadJSON(providerConfigPath, config)
		if(err != nil) {
			fmt.Println(err)
			return nil
		}
		providerConfigCache[providerName] = config
		return config
	} else {
		return providerConfigCache[providerName]
	}
}

func GetPackageConfig() (utils.Json, error) {
	var file = utils.FileExists(ProviderPath + "/GetPackageConfig.js")
	if(file) {
		input := make(map[string]interface {})		
		res, err :=  jsVm.Run(ProviderPath + "/GetPackageConfig.js", input)
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
	var file = utils.FileExists(ProviderPath + "/PostGetPackageConfig.js")
	if(file) {
		input := make(map[string]interface {})
		input["PackageConfig"] = config
		
		res, err :=  jsVm.Run(ProviderPath + "/PostGetPackageConfig.js", input)
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
	var file = utils.FileExists(ProviderPath + "/SaveDependencyList.js")
	if(file) {
		input := make(map[string]interface {})
		input["Dependencies"] = depList
		
		_, err :=  jsVm.Run(ProviderPath + "/SaveDependencyList.js", input)
		if(err != nil) {
			return err
		}

		return nil
	} else {
		return defaultProvider.SaveDependencyList(depList)
	}
}

func GetDependencyList(config utils.Json) ([]map[string]interface {}, error) {
	var file = utils.FileExists(ProviderPath + "/GetDependencyList.js")
	if(file) {
		input := make(map[string]interface {})
		input["PackageConfig"] = config
		
		res, err :=  jsVm.Run(ProviderPath + "/GetDependencyList.js", input)
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
	depProviderPath := utils.DIRNAME() + "/plugins/provider-" + dependency["provider"].(string)
	var file = utils.FileExists(depProviderPath + "/ResolveDependencyLocation.js")
	if(dependency["provider"].(string) != "gupm" && file) {
		input := make(map[string]interface {})
		input["Dependency"] = dependency

		res, err :=  jsVm.Run(depProviderPath + "/ResolveDependencyLocation.js", input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		
		if(resObj == nil) {
			fmt.Println("ERROR Failed to resolve", dependency, "Trying again.")
			return ResolveDependencyLocation(dependency)
		}
		return resObj.(map[string]interface {}), err1
	} else {
		return nil, nil
	}
}

func ExpandDependency(dependency map[string]interface {}) (map[string]interface {}, error) {
	depProviderPath := utils.DIRNAME() + "/plugins/provider-" + dependency["provider"].(string)
	var file = utils.FileExists(depProviderPath + "/ExpandDependency.js")
	if(dependency["provider"].(string) != "gupm" && file) {
		input := make(map[string]interface {})
		input["Dependency"] = dependency

		res, err :=  jsVm.Run(depProviderPath + "/ExpandDependency.js", input)
		if(err != nil) {
			return nil, err
		}

		toExport, _ := res.Export()
		resObj := utils.JsonExport(toExport).(map[string] interface {})
		// fmt.Println(resObj["dependencies"])

		if(resObj == nil) {
			fmt.Println("ERROR Failed to resolve", dependency, ". Trying again.")
			return ExpandDependency(dependency)
		}

		return resObj, nil
	} else {
		return dependency, nil
	}
}

func GetDependency(provider string, name string, version string, url string, path string) (string, error) {
	depProviderPath := utils.DIRNAME() + "/plugins/provider-" + provider
	var file = utils.FileExists(depProviderPath + "/GetDependency.js")
	if(provider != "gupm" && file) {
		input := make(map[string]interface {})
		input["Provider"] = provider
		input["Name"] = name
		input["Version"] = version
		input["Url"] = url
		input["Path"] = path

		res, err :=  jsVm.Run(depProviderPath + "/GetDependency.js", input)
		if(err != nil) {
			return "", err
		}

		resStr, err1 := res.ToString()
		return resStr, err1
	} else {
		return defaultProvider.GetDependency(provider, name, version, url, path)
	}
}

func BinaryInstall(path string) (error) {
	os.RemoveAll("./.bin")
	var file = utils.FileExists(ProviderPath + "/BinaryInstall.js")
	if(file) {
		input := make(map[string]interface {})
		input["Destination"] = ".bin"
		input["Source"] = "node_modules"

		res, err :=  jsVm.Run(ProviderPath + "/BinaryInstall.js", input)
		if(err != nil) {
			return err
		}

		_, err1 := res.ToString()
		return err1
	} else {
		return nil
	}
}

func PostGetDependency(provider string, name string, version string, url string, path string, result string) (string, error) {
	depProviderPath := utils.DIRNAME() + "/plugins/provider-" + provider
	var file = utils.FileExists(depProviderPath + "/PostGetDependency.js")
	if(provider != "gupm" && file) {
		input := make(map[string]interface {})
		input["Provider"] = provider
		input["Name"] = name
		input["Version"] = version
		input["Url"] = url
		input["Path"] = path
		input["Result"] = result

		res, err :=  jsVm.Run(depProviderPath + "/PostGetDependency.js", input)
		if(err != nil) {
			return "", err
		}

		resStr, err1 := res.ToString()
		return resStr, err1
	} else {
		return defaultProvider.PostGetDependency(provider, name, version, url, path, result)
	}
}

func eliminateRedundancy(tree []map[string]interface {}, path map[string]bool) []map[string]interface {} {
	var cleanTree = make([]map[string]interface {}, 0)
	for index, dep := range tree {
		if(dep["name"] != nil) {
			_ = index
			depKey := dep["name"].(string) + "@" + dep["version"].(string)
			if(path[depKey] != true) {
				cleanTree = append(cleanTree, dep)
			}
		}
	}
	
	for index, dep := range cleanTree {
		if(dep["name"] != nil) {
			nextDepList, ok := dep["dependencies"].([]map[string]interface {})

			if(ok) {
				depKey := dep["name"].(string) + "@" + dep["version"].(string)
				newPath := make(map[string]bool)
				for key, value := range path {
					newPath[key] = value
				}
				newPath[depKey] = true
				newSubTree := eliminateRedundancy(nextDepList, newPath)
				cleanTree[index]["dependencies"] = newSubTree
			}
		}
	}
	return cleanTree
}

func flattenDependencyTree(tree []map[string]interface {}, subTree []map[string]interface {}) ([]map[string]interface {}, []map[string]interface {}) {
	var cleanTree = make([]map[string]interface {}, 0)

	for index, dep := range subTree {
		var rootDeps = make(map[string]string)

		for _, dep := range tree {
			rootDeps[dep["name"].(string)] = dep["version"].(string)
		}

		if(rootDeps[dep["name"].(string)] == "") {
			tree = append(tree, dep)

			nextDepList, ok := dep["dependencies"].([]map[string]interface {})
	
			if(ok) {
				newTree, newSubTree := flattenDependencyTree(tree, nextDepList)
				tree = newTree
				subTree[index]["dependencies"] = newSubTree
			}
		} else if(rootDeps[dep["name"].(string)] != dep["version"].(string)) {
			nextDepList, ok := dep["dependencies"].([]map[string]interface {})
	
			if(ok) {
				newTree, newSubTree := flattenDependencyTree(tree, nextDepList)
				tree = newTree
				subTree[index]["dependencies"] = newSubTree
			}

			cleanTree = append(cleanTree, subTree[index])
		}
	}

	return tree, cleanTree
}

func BuildDependencyTree(tree []map[string]interface {}) []map[string]interface {} {
	cleanTree := eliminateRedundancy(tree, make(map[string]bool))

	for index, dep := range cleanTree {
		nextDepList, ok := dep["dependencies"].([]map[string]interface {})

		if(ok) {
			newCleanTree, newDepList := flattenDependencyTree(cleanTree, nextDepList)
			cleanTree = newCleanTree
			cleanTree[index]["dependencies"] = newDepList
		}
	}
	return cleanTree
}

func installDependencySubFolders(path string, depPath string) {
	files := utils.ReadDir(path)

	for _, file := range files {
		if(file.IsDir()) {
			folderPath := depPath + "/" + file.Name()
			os.MkdirAll(folderPath, os.ModePerm);
			installDependencySubFolders(path + "/" + file.Name(), folderPath)
		} else {
			os.Link(path + "/" + file.Name(), depPath + "/" + file.Name())
		}
	}
}

func InstallDependency(path string, dep map[string]interface {}) {
	depPath := path + "/" + dep["name"].(string)
	// if(utils.FileExists(depPath)) {
	// 	// TODO: check version
	// } else {
	// }
	os.MkdirAll(depPath, os.ModePerm);
	installDependencySubFolders(dep["path"].(string), depPath)
}