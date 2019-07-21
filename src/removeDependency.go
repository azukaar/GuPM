package main

import (
	"./provider"
	"./ui"
	"./utils"
)

func remove(slice []map[string]interface{}, s int) []map[string]interface{} {
	return append(slice[:s], slice[s+1:]...)
}

func RemoveDependency(path string, rls []string) error {
	var err error
	var packageConfig utils.Json
	var depList []map[string]interface{}

	ui.Title("Add dependency...")

	err = provider.InitProvider(Provider)

	if err != nil {
		return err
	}

	providerConfig, err = provider.GetProviderConfig(Provider)
	ui.Error(err)
	packageConfig, _ = provider.GetPackageConfig(path)
	packageConfig, _ = provider.PostGetPackageConfig(packageConfig)

	depList, err = provider.GetDependencyList(packageConfig)
	if err != nil {
		return err
	}

	ui.Title("Removing from dependency list...")

	for _, str := range rls {
		for index, dep := range depList {
			if dep["name"].(string) == str {
				depList = remove(depList, index)
			}
		}
	}

	err = provider.SaveDependencyList(path, depList)
	if err != nil {
		return err
	}

	// TODO: Remove from module folder

	return nil
}
