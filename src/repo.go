package main

import (
	"./utils"
	"io/ioutil"
)

func GetOrCreateRepo(path string) map[string]interface{} {
	configPath := path + "/gupm_repo.json"
	if(!utils.FileExists(configPath)) {
		baseConfig := `{
	"packages": {}
}`
		utils.WriteFile(configPath, baseConfig)
		return utils.StringToJSON(baseConfig)
	}

	file, _ := ioutil.ReadFile(configPath)
	return utils.StringToJSON(string(file))
}

func SaveRepo(path string, file map[string]interface{}) {
	utils.WriteJsonFile(path + "/gupm_repo.json", file)
}