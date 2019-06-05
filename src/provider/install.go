package provider

import (
	"../utils"
	"../jsVm"
	"../ui"
	"os"
	"../defaultProvider"
	"regexp"
	"io/ioutil"
)

func BinaryInstall(path map[string]string) (error) {
	os.RemoveAll(".bin")

	for pr, prdir := range path {
		depProviderPath := utils.Path(utils.DIRNAME() + utils.Path("/plugins/provider-" + pr))
		var file = utils.FileExists(depProviderPath + utils.Path("/binaryInstall.gs"))
		if(pr != "gupm" && file) {
			input := make(map[string]interface {})
			input["Destination"] = ".bin"
			input["Source"] = prdir
	
			res, err :=  jsVm.Run(depProviderPath + utils.Path("/binaryInstall.gs"), input)
			if(err != nil) {
				return err
			}
	
			_, err1 := res.ToString()
			return err1
		} else {
			return defaultProvider.BinaryInstall(".bin", prdir)
		}
	}
	return nil
}

func installDependencySubFolders(path string, depPath string) {
	files := utils.ReadDir(path)

	for _, file := range files {
		if(file.IsDir()) {
			folderPath := utils.Path(depPath + "/" + file.Name())
			os.MkdirAll(folderPath, os.ModePerm);
			installDependencySubFolders(utils.Path(path + "/" + file.Name()), folderPath)
		} else {
			isFileExists := false
			err := os.Link(utils.Path(path + "/" + file.Name()), utils.Path(depPath + "/" + file.Name()))
			if(err != nil) {
				isFileExists, _ = regexp.MatchString(`file exists$`, err.Error())
			}

			if(err != nil && !isFileExists) {
				if(!linkHasErrored) {
					ui.Error(err.Error())
					ui.Error("Error, cannot use hard link on your system. Falling back to copying file (Will be slower!)")
					linkHasErrored = true
				}
				input, err := ioutil.ReadFile(utils.Path(path + "/" + file.Name()))
        if err != nil {
                ui.Error(err.Error())
                return
        }

        err = ioutil.WriteFile(utils.Path(depPath + "/" + file.Name()), input, 0644)
        if err != nil {
                ui.Error(err.Error())
                return
        }
			}
		}
	}
}

func InstallDependency(path string, dep map[string]interface {}) {
	depPath := utils.Path(path + "/" + dep["name"].(string))
	// if(utils.FileExists(depPath)) {
	// 	// TODO: check version
	// } else {
	// }
	_, ok := dep["path"].(string)
	if(ok) {
		os.MkdirAll(utils.Path(depPath), os.ModePerm);
		installDependencySubFolders(utils.Path(dep["path"].(string)), depPath)
	} else {
		ui.Error(dep["name"].(string) + " Cannot be installed.")
	}
}