package main

import (
	"./provider"
	"os"
	"./utils"
	"fmt"
)

var cacheExpanded = make(map[string]map[string]interface {})

func expandDepList(depList []map[string]interface {}) []map[string]interface {} {
	var channel = make(chan int)
	var todo = len(depList)
	for index, dep := range depList {
		go (func(channel chan int, index int, dep map[string]interface {}){
			if(dep["expanded"] != true) {
				newDep, errExpand := provider.ResolveDependencyLocation(dep)
				if(newDep == nil) {
					fmt.Println("Error: Provider NPM didnt resolve ", dep)
				}
				newDep["path"] = utils.DIRNAME() + "/cache/" + newDep["provider"].(string) + "/" + newDep["name"].(string) + "/" + newDep["version"].(string)

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
			channel <- 1
		})(channel, index, dep)
	}

	for _,_ = range depList {
		<-channel
	}

	return depList
}

func installDep(path string, depList []map[string]interface {}) {
	var channel = make(chan int)
	fmt.Println("Installing...", path)
	for index, dep := range depList {
		go (func(channel chan int, index int, dep map[string]interface {}){
			provider.InstallDependency(path, dep)

			nextDepList, ok := depList[index]["dependencies"].([]map[string]interface {})

			if(ok) {
				installDep(path + "/" + depList[index]["name"].(string) + "/" + providerConfig.Config.Default.InstallPath, nextDepList)
			}
			channel <- 1
		})(channel, index, dep)
	}
	for _,_ = range depList {
		<-channel
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
	fmt.Println("Expand dependency list...")
	depList = expandDepList(depList)
	
	fmt.Println("Build dependency list...")
	depList = provider.BuildDependencyTree(depList)
	
	os.MkdirAll(providerConfig.Config.Default.InstallPath, os.ModePerm);
	
	fmt.Println("Install dependencies...")
	installDep(providerConfig.Config.Default.InstallPath, depList)

	fmt.Println("Install Binaries...")
	err = provider.BinaryInstall(providerConfig.Config.Default.InstallPath)
	if(err != nil) {
		return err
	}
	
	return nil
}