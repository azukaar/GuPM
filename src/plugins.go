package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"./provider"
	"./ui"
	"./utils"
)

func PluginLink(path string) {
	configPath := utils.Path(path + "/gupm.json")
	if utils.FileExists(configPath) {
		packageConfig := new(utils.GupmEntryPoint)
		errConfig := utils.ReadJSON(configPath, &packageConfig)
		if errConfig != nil {
			ui.Error("Can't read provider configuration")
			ui.Error(errConfig)
			return
		}

		pluginFolder := utils.HOMEDIR(".") + utils.Path("/.gupm/plugins/")
		os.MkdirAll(pluginFolder, os.ModePerm)
		err := os.Symlink(utils.AbsPath(path), pluginFolder+packageConfig.Name)
		if err != nil {
			ui.Error(err)
		}
	} else {
		ui.Error("Can't find provider configuration")
	}
}

func PluginInstall(path string, plugins []string) error {
	ui.Title("Install plugin...")
	pluginFolder := utils.HOMEDIR(".") + utils.Path("/.gupm/plugins/")

	for _, rls := range plugins {
		ui.Log(rls)
		newDep, err := provider.ResolveDependencyLocation(utils.BuildDependencyFromString("https", rls))
		if err != nil {
			return err
		}

		pluginName := filepath.Base(newDep["name"].(string))
		if len(strings.Split(pluginName, ":")) > 1 {
			pluginName = strings.Split(pluginName, ":")[1]
		}
		if len(strings.Split(pluginName, "/")) > 1 {
			pluginName = strings.Split(pluginName, "/")[1]
		}

		newDep["path"] = pluginFolder + utils.Path("/"+pluginName)
		getRes, errorGD := provider.GetDependency(
			newDep["provider"].(string),
			newDep["name"].(string),
			newDep["version"].(string),
			newDep["url"].(string),
			newDep["path"].(string),
		)
		if errorGD != nil {
			return errorGD
		}
		_, errorPGD := provider.PostGetDependency(
			newDep["provider"].(string),
			newDep["name"].(string),
			newDep["version"].(string),
			newDep["url"].(string),
			newDep["path"].(string),
			getRes,
		)
		if errorPGD != nil {
			return errorPGD
		}
	}
	return nil
}

func PluginDelete(path string, plugins []string) {
	folders := make([]string, 0)
	pluginFolder := utils.HOMEDIR(".") + utils.Path("/.gupm/plugins/")

	for _, str := range plugins {
		folders = append(folders, pluginFolder+str)
	}

	utils.RemoveFiles(folders)
	fmt.Println("Done deleting.")
}

func PluginCreate(path string) {
	fmt.Println("Welcome to the plugin creation assistant")
	name := "provider-" + ui.WaitForInput("What is the name of the plugin? provider-")
	description := ui.WaitForInput("Enter a description: ")
	author := ui.WaitForInput("Enter the author: ")
	licence := ui.WaitForInput("Enter the licence (ISC): ")
	ppath := utils.Path(path + "/" + name)

	os.MkdirAll(ppath, os.ModePerm)
	os.MkdirAll(ppath+utils.Path("/docs/repo"), os.ModePerm)

	utils.WriteFile(ppath+utils.Path("/gupm.json"), `{
	"name": "`+name+`",
	"version": "0.0.1",
	"description": "`+description+`",
	"author": "`+author+`",
	"licence": "`+licence+`",
    "publish": {
        "source": ["."],
        "dest": "../docs/repo"
    },
    "config": {
        "default": {
            "entrypoint": "gupm.json",
            "installPath": "gupm_modules"
        }
    }
}`)

	fmt.Println("creation done.")
}
