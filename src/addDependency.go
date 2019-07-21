package main

import (
	"./provider"
	"./ui"
	"./utils"
)

func AddDependency(path string, rls []string) error {
	var err error
	var packageConfig utils.Json
	var depList []map[string]interface{}
	var depProvider = Provider

	ui.Title("Add dependency...")

	if !ProviderWasForced && utils.FileExists(path+utils.Path("/gupm.json")) {
		config, _ := utils.ReadGupmJson(path + utils.Path("/gupm.json"))
		if config.Dependencies.DefaultProvider != "" {
			depProvider = config.Dependencies.DefaultProvider
		}
	}

	err = provider.InitProvider(Provider)

	if err != nil {
		return err
	}

	providerConfig, err = provider.GetProviderConfig(Provider)
	if err != nil {
		return err
	}

	packageConfig, err = provider.GetPackageConfig(path)
	if err != nil {
		return err
	}

	packageConfig, err = provider.PostGetPackageConfig(packageConfig)
	if err != nil {
		return err
	}

	depList, err = provider.GetDependencyList(packageConfig)
	if err != nil {
		return err
	}

	ui.Title("Adding to dependency list...")

	for _, str := range rls {
		dep := utils.BuildDependencyFromString(depProvider, str)
		resolved, err := provider.ResolveDependencyLocation(dep)
		if err != nil || resolved["url"].(string) == "" {
			ui.Error("Can't resolve", str)
			return err
		}
		dep["version"] = resolved["version"]
		depList = append(depList, dep)
	}

	if packageConfig != nil {
		err = provider.SaveDependencyList(path, depList)
		if err != nil {
			return err
		}
	}

	return nil
}
