package main

import (
	"./provider"
	"./utils"
	"./ui"
)

func AddDependency(path string, rls []string) error {	
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

	ui.Title("Adding to dependency list...")

	for _, str := range rls {
		dep :=  utils.BuildDependencyFromString(Provider, str)
		resolved, err := provider.ResolveDependencyLocation(dep)
		if(err != nil) {
			ui.Error("Error can't resolve : " + dep["Name"].(string) + "@" +  dep["version"].(string))
			return err
		}
		dep["version"] = resolved["version"]
		depList = append(depList, dep)
	}

	err = provider.SaveDependencyList(depList)
	if(err != nil) {
		return err
	}

	InstallProject(".")

	return nil
}