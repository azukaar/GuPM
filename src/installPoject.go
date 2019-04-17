package main

import (
	"./provider"
	"./utils"
	"fmt"
)

var cacheExpanded = make(map[string]map[string]interface {})

func expandDepList(depList []map[string]interface {}) []map[string]interface {} {
	var todo = len(depList)
	for index, dep := range depList {
		if(dep["expanded"] != true) {
			newDep, errExpand := provider.ResolveDependencyLocation(dep)
			newDep["path"] = "./cache/" + newDep["provider"].(string) + "/" + newDep["name"].(string) + "/" + newDep["version"].(string)

			if(!utils.FileExists(newDep["path"].(string))) {
				getRes, errorGD := provider.GetDependency(
					newDep["provider"].(string),
					newDep["name"].(string),
					newDep["version"].(string),
					newDep["url"].(string),
					newDep["path"].(string),
				)
				if(errorGD != nil) {
					fmt.Println(1, errorGD)
				}
				_, errorPGD := provider.PostGetDependency(
					newDep["provider"].(string),
					newDep["name"].(string),
					newDep["version"].(string),
					newDep["url"].(string),
					newDep["path"].(string),
					getRes,
				)
				if(errorPGD != nil) {
					fmt.Println(2, errorPGD)
				}
			}

			if(newDep["expanded"] != true) {
				if(cacheExpanded[newDep["url"].(string)]["expanded"] != true) {
					newDep, errExpand = provider.ExpandDependency(newDep)
					if(errExpand != nil) {
						fmt.Println(errExpand)
					}
	
					newDep["expanded"] = true
					cacheExpanded[newDep["url"].(string)] = newDep
				} else {
					newDep = cacheExpanded[newDep["url"].(string)]
				}
			}

			depList[index] = newDep

			fmt.Println("Get dependency", index, "of", todo)
			
			nextDepList, ok := depList[index]["dependencies"].([]map[string]interface {})

			if(ok) {
				depList[index]["dependencies"] = expandDepList(nextDepList)
			}
		}
	}

	return depList
}

func installDep(path string, depList []map[string]interface {}) {
	fmt.Println("Installing...", path)
	for index, dep := range depList {
		
		provider.InstallDependency(path, dep)

		nextDepList, ok := depList[index]["dependencies"].([]map[string]interface {})

		if(ok) {
			installDep(path + "/" + depList[index]["name"].(string) + "/" + providerConfig.Config.Default.InstallPath, nextDepList)
		}
	}
}

var providerConfig *provider.GupmEntryPoint

func InstallProject(path string) error {
	_ = fmt.Println
	
	var err error
	var packageConfig utils.Json
	var depList []map[string]interface {}

	err = provider.InitProvider(Provider)
	if(err != nil) {
		return err
	}

	providerConfig = provider.GetProviderConfig()
	packageConfig, _ = provider.GetPackageConfig()
	packageConfig, _ = provider.PostGetPackageConfig(packageConfig)

	depList, err = provider.GetDependencyList(packageConfig)
	if(err != nil) {
		return err
	}

	depList = expandDepList(depList)

	depList = provider.BuildDependencyTree(depList)

	installDep(providerConfig.Config.Default.InstallPath, depList)
	
	// fmt.Println(depList)
	
	// depList.foreach
	// provider.postInstallation(path, packageConfig)
	
	// provider.finalHook(path, packageConfig)

	return nil
}