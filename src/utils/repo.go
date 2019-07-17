package utils

import (
	"io/ioutil"
)

func GetOrCreateRepo(path string) map[string]interface{} {
	configPath := path + "/gupm_repo.json"
	if !FileExists(configPath) {
		baseConfig := `{
	"packages": {}
}`
		WriteFile(configPath, baseConfig)
		return StringToJSON(baseConfig)
	}

	file, _ := ioutil.ReadFile(configPath)
	return StringToJSON(string(file))
}

func SaveRepo(path string, file map[string]interface{}) {
	WriteJsonFile(path+"/gupm_repo.json", file)
}
