package main

import (
	"./provider"
	"./utils"
	"fmt"
)

func InstallProject(path string) error {
	var err error
	var packageConfig utils.Json
	var depList utils.PackageDepedencyListType

	if(Provider != "") {
		err = provider.InitProvider(Provider)
		if(err != nil) {
			return err
		}
	}

	packageConfig, _ = provider.GetPackageConfig()
	packageConfig, _ = provider.PostGetPackageConfig(packageConfig)

	depList, err = provider.GetDependencyList(packageConfig)
	if(err != nil) {
		return err
	}
	
	for index, dep := range depList {
		expDep, errExpand := provider.ExpandDependency(dep.(utils.Json))
		if(errExpand != nil) {
			return errExpand
		}
		depList[index] = expDep
	}

	// depList = provider.BuildDependencyTree(depList)
	fmt.Println(depList)
	
	// depList.foreach
	// provider.getDepedency()
	// provider.postInstallation(path, packageConfig)
	
	// provider.finalHook(path, packageConfig)

	return nil
}