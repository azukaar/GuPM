package defaultProvider

import (
	"../utils"
	"os"
	"../ui"
)

func Publish(path string) error {
	configPath := utils.Path(path + "/gupm.json")
	if(utils.FileExists(configPath)) {
		packageConfig := new(utils.GupmEntryPoint)
		errConfig := utils.ReadJSON(configPath, &packageConfig)
		if(errConfig != nil) {
			ui.Error("Can't read provider configuration")
			return errConfig
		}
	
		ppath := utils.Path(path + "/" + packageConfig.Publish.Dest)
		repoConfig := utils.GetOrCreateRepo(ppath)
		packageList := repoConfig["packages"].(map[string]interface{})

		if(packageList[packageConfig.Name] != nil) {
			if(utils.Contains(packageList[packageConfig.Name], packageConfig.Version)) {
				ui.Error("Package " + packageConfig.Name + "@" + packageConfig.Version + " already published. Please bump the version number.")
				return nil
			} else {
				packageList[packageConfig.Name] = append(utils.ArrString(packageList[packageConfig.Name]), packageConfig.Version)
			}
		} else {
			packageList[packageConfig.Name] = make([]string, 0)
			packageList[packageConfig.Name] = append(utils.ArrString(packageList[packageConfig.Name]), packageConfig.Version)
		}

		installPath := ppath + utils.Path("/" + packageConfig.Name + "/" + packageConfig.Version)
		os.MkdirAll(installPath, os.ModePerm)

		sourcePaths := make([]string, 0)
		for _, src := range packageConfig.Publish.Source {
			sourcePaths = append(sourcePaths, utils.Path(path + "/" + src));
		}
		arch, _ := utils.Tar(sourcePaths)
		arch.SaveAt(installPath + utils.Path("/" + packageConfig.Name + "-" + packageConfig.Version + ".tgz"))

		repoConfig["packages"] = packageList
		utils.SaveRepo(ppath, repoConfig)
	} else {
		ui.Error("Can't find provider configuration")
	}

	return nil
}