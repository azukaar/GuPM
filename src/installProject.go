package main

import (
	"./provider"
	"os"
	"errors"
	"./utils"
	"./ui"
	"sync"
)

var cacheExpanded = make(map[string]map[string]interface {})
var lock = sync.RWMutex{}

func expandDepList(depList []map[string]interface {}) ([]map[string]interface {}) {
	channel := make(chan int)
	// todo := len(depList)
	
	for index, dep := range depList {
		go (func(channel chan int, index int, dep map[string]interface {}) {
			if(dep["expanded"] != true) {
				newDep := make(map[string]interface {})
				for key, value := range dep {
					newDep[key] = value
				}
				newDep, errExpand := provider.ResolveDependencyLocation(newDep)
				if(newDep == nil) {
					ui.Error("Error: Provider " + dep["provider"].(string) + " didnt resolve " + dep["name"].(string) + "@" + dep["version"].(string))
					ui.Error(errExpand.Error())
					channel <- 0
					return;
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
						ui.Error(errorGD.Error())
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
						ui.Error(errorPGD.Error())
					}
				}

				lock.Lock()
				if(newDep["expanded"] != true) {
					if(cacheExpanded[newDep["url"].(string)]["expanded"] != true) {
						newDep, errExpand = provider.ExpandDependency(newDep)
						if(errExpand != nil || newDep == nil || len(newDep) == 0)  {
							ui.Error("Error: Provider " + dep["provider"].(string) + " didnt expand " + dep["name"].(string) + "@" + dep["version"].(string))
							ui.Error(errExpand.Error())
							channel <- 0
							return;
						}
		
						newDep["expanded"] = true
						cacheExpanded[newDep["url"].(string)] = newDep
					} else {
						newDep = cacheExpanded[newDep["url"].(string)]
					}
				}
				lock.Unlock()

				depList[index] = newDep

				ui.Log("Get dependency " + newDep["name"].(string))
				
				nextDepList, ok := depList[index]["dependencies"].([]map[string]interface {})

				if(ok) {
					depList[index]["dependencies"] = expandDepList(nextDepList)
				}
			}
			channel <- 1
		})(channel, index, dep)
	}

	for _,_ = range depList {
		i := <-channel
		if(i == 0) {
			return nil
		}
	}

	return depList
}

func installDep(path string, depList []map[string]interface {}) {
	var channel = make(chan int)
	ui.Log("Installing " + path)
	for index, dep := range depList {
		go (func(channel chan int, index int, dep map[string]interface {}){
			depProviderConfig := provider.GetProviderConfig(dep["provider"].(string))
			provider.InstallDependency(path, dep)

			nextDepList, ok := depList[index]["dependencies"].([]map[string]interface {})

			if(ok) {
				installDep(path + "/" + depList[index]["name"].(string) + "/" + depProviderConfig.Config.Default.InstallPath, nextDepList)
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
	ui.Title("Installing project...")
	
	var err error
	var packageConfig utils.Json
	var depList []map[string]interface {}

	err = provider.InitProvider(Provider)
	if(err != nil) {
		return err
	}

	providerConfig = provider.GetProviderConfig(Provider)
	packageConfig, _ = provider.GetPackageConfig()
	packageConfig, _ = provider.PostGetPackageConfig(packageConfig)

	depList, err = provider.GetDependencyList(packageConfig)
	if(err != nil) {
		return err
	}
	if(depList == nil) {
		return errors.New("Failed to get dependancy list")
	}

	ui.Title("Expand dependency list...")
	depList = expandDepList(depList)
	if(depList == nil) {
		return errors.New("Failed to expand dependancy list")
	}
	
	ui.Title("Build dependency list...")
	depList = provider.BuildDependencyTree(depList)
	if(depList == nil) {
		return errors.New("Failed to build dependancy list")
	}
	
	os.MkdirAll(providerConfig.Config.Default.InstallPath, os.ModePerm);
	
	ui.Title("Install dependencies...")
	installDep(providerConfig.Config.Default.InstallPath, depList)

	ui.Title("Install Binaries...")
	err = provider.BinaryInstall(providerConfig.Config.Default.InstallPath)
	if(err != nil) {
		return err
	}

	ui.Title("Installation done ❤️")
	
	return nil
}