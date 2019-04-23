package provider

import (
	"fmt"
	"../defaultProvider"
	"../utils"
	"../jsVm"
	"io/ioutil"
	"os"
	"github.com/robertkrimen/otto"
)

var Provider string
var ProviderPath string
var ProviderConfig *GupmEntryPoint
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
	if(Provider != "") {
		ProviderPath = "plugins/provider-" + Provider
		fmt.Println("Reading provider config for", Provider)
		ProviderConfig = new(GupmEntryPoint) 
		err := utils.ReadJSON(ProviderPath + "/gupm.json", ProviderConfig)
		if(err != nil) {
			return err
		} else {
			fmt.Println("Initialisation OK for", ProviderConfig.Name);
		}
	} else {
		ProviderConfig = new(GupmEntryPoint) 
		err := utils.ReadJSON("gupm.json", ProviderConfig)
		if(err != nil) {
			return err
		} else {
			fmt.Println("Initialisation OK for", ProviderConfig.Name);
		}
	}
	return nil
}

func GetEntryPoint() string {
	return ProviderConfig.Config.Default.Entrypoint
}

func GetProviderConfig() *GupmEntryPoint {
	if(Provider != "") {
		return ProviderConfig
	} else {
		return new(GupmEntryPoint) 
	}
}

func GetPackageConfig() (utils.Json, error) {
	var file = utils.FileExists(ProviderPath + "/GetPackageConfig.js")
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
	var file = utils.FileExists(ProviderPath + "/PostGetPackageConfig.js")
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
	var file = utils.FileExists(ProviderPath + "/GetDependencyList.js")
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

func ResolveDependencyLocation(dependency map[string]interface {}) (map[string]interface {}, error) {
	var file = utils.FileExists(ProviderPath + "/ResolveDependencyLocation.js")
	if(file) {
		input := make(map[string]interface {})
		input["Dependency"] = dependency

		res, err :=  run(ProviderPath + "/ResolveDependencyLocation.js", input)
		if(err != nil) {
			return nil, err
		}

		resObj, err1 := res.Export()
		return resObj.(map[string]interface {}), err1
	} else {
		return nil, nil
	}
}

func ExpandDependency(dependency map[string]interface {}) (map[string]interface {}, error) {
	var file = utils.FileExists(ProviderPath + "/ExpandDependency.js")
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
	var file = utils.FileExists(ProviderPath + "/GetDependency.js")
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
	var file = utils.FileExists(ProviderPath + "/PostGetDependency.js")
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
	var rootDeps = make(map[string]string)
	var cleanTree = make([]map[string]interface {}, 0)

	for _, dep := range tree {
		rootDeps[dep["name"].(string)] = dep["version"].(string)
	}
	
	for index, dep := range subTree {
		if(rootDeps[dep["name"].(string)] == dep["version"].(string)) {
			// delete
		} else {
			if (rootDeps[dep["name"].(string)] == "") {
				tree = append(tree, dep)
			}

			nextDepList, ok := dep["dependencies"].([]map[string]interface {})
	
			if(ok) {
				newTree, newSubTree := flattenDependencyTree(tree, nextDepList)
				tree = newTree
				subTree[index]["dependencies"] = newSubTree
			}

			if (rootDeps[dep["name"].(string)] != "") {
				cleanTree = append(cleanTree, subTree[index])
			}
		}
	}

	return tree, cleanTree
}

func BuildDependencyTree(tree []map[string]interface {}) []map[string]interface {} {
	cleanTree := eliminateRedundancy(tree, make(map[string]bool))
	for index, dep := range cleanTree {
		nextDepList, ok := dep["dependencies"].([]map[string]interface {})

		if(ok) {
			newCleanTree, newDepList := flattenDependencyTree(tree, nextDepList)
			cleanTree = newCleanTree
			cleanTree[index]["dependencies"] = newDepList
		}
	}
	return cleanTree
}

func readDir(path string) []os.FileInfo{
    files, err := ioutil.ReadDir(path)
    if err != nil {
        fmt.Println(err)
	}

    return files
}

func InstallDependency(path string, dep map[string]interface {}) {
	completePath := path + "/" + dep["name"].(string)
	
	_, ok := dep["dependencies"].([]map[string]interface {})
	
	if(ok && len(dep["dependencies"].([]map[string]interface {})) > 0) {
		os.MkdirAll(completePath, os.ModePerm);
		files := readDir(dep["path"].(string))
		for _, file := range files {
			os.Symlink(dep["path"].(string) + "/" + file.Name(), completePath + "/" + file.Name())
		}
	} else {
		os.Symlink(dep["path"].(string), completePath)
	}
}