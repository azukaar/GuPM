package main

import (
	"./provider"
	"./utils"
	"fmt"
)

func AddDependency(path string, rls []string) error {
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

	fmt.Println("Adding to dependency list...")

	for _, str := range rls {
		dep :=  utils.BuildDependencyFromString(Provider, str)
		resolved, err := provider.ResolveDependencyLocation(dep)
		if(err != nil) {
			fmt.Println(1, dep)
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