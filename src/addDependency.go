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

	ui.Title("Add dependency...")

	err = provider.InitProvider(Provider)

	if(err != nil) {
		return err
	}

	providerConfig, err = provider.GetProviderConfig(Provider)
	if(err != nil) {
		return err
	}
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
			ui.Error("Can't resolve", str)
			return err
		}
		dep["version"] = resolved["version"]
		depList = append(depList, dep)
	}

	err = provider.SaveDependencyList(depList)
	if(err != nil) {
		return err
	}

	return nil
}