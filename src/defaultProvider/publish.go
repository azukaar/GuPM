package defaultProvider

import (
	"os"

	"../ui"
	"../utils"
)

func Publish(path string, namespace string) error {
	configPath := utils.Path(path + "/gupm.json")
	if utils.FileExists(configPath) {
		packageConfig := new(utils.GupmEntryPoint)
		errConfig := utils.ReadJSON(configPath, &packageConfig)
		if errConfig != nil {
			ui.Error("Can't read provider configuration")
			return errConfig
		}

		pname := packageConfig.Name

		if namespace != "" {
			pname = namespace + "/" + pname
		}

		ppath := utils.Path(path + "/" + packageConfig.Publish.Dest)
		repoConfig := utils.GetOrCreateRepo(ppath)
		packageList := repoConfig["packages"].(map[string]interface{})

		if packageList[pname] != nil {
			if utils.Contains(packageList[pname], packageConfig.Version) {
				ui.Error("Package " + pname + "@" + packageConfig.Version + " already published. Please bump the version number.")
				return nil
			} else {
				packageList[pname] = append(utils.ArrString(packageList[pname]), packageConfig.Version)
			}
		} else {
			packageList[pname] = make([]string, 0)
			packageList[pname] = append(utils.ArrString(packageList[pname]), packageConfig.Version)
		}

		installPath := ppath + utils.Path("/"+pname+"/"+packageConfig.Version)
		os.MkdirAll(installPath, os.ModePerm)

		sourcePaths := make([]string, 0)
		for _, src := range packageConfig.Publish.Source {
			sourcePaths = append(sourcePaths, utils.Path(path+"/"+src))
		}
		arch, _ := utils.Tar(sourcePaths)
		arch.SaveAt(installPath + utils.Path("/"+packageConfig.Name+"-"+packageConfig.Version+".tgz"))

		repoConfig["packages"] = packageList
		utils.SaveRepo(ppath, repoConfig)
	} else {
		ui.Error("Can't find provider configuration")
	}

	return nil
}
