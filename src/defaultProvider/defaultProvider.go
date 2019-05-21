package defaultProvider

import (
	"io/ioutil"
	"encoding/json"
	"../utils"
	"../ui"
	"os"
	"regexp"
	"reflect"
)

func SaveDependencyList(depList []map[string]interface{}) error {
	_ = depList
	return nil
}

func Bootstrap(path string) {
	
}

func GetPackageConfig(entryPoint string) map[string]interface {} {
	var packageConfig map[string]interface{}
	b, err := ioutil.ReadFile(entryPoint)
	if err != nil {
		ui.Error(err.Error() + " : " + entryPoint)
	}

	json.Unmarshal([]byte(string(b)), &packageConfig)

	return packageConfig
}

func GetDependency(provider string, name string, version string, url string, path string) (string, error) {
	return string(utils.HttpGet(url)), nil
}

func PostGetDependency(provider string, name string, version string, url string, path string, result string) (string, error) {
	os.MkdirAll(path, os.ModePerm)
	tarCheck := regexp.MustCompile(`\.tgz$`)
	tryTar := tarCheck.FindString(url)
	zipCheck := regexp.MustCompile(`\.zip$`)
	tryZip := zipCheck.FindString(url)

	if(tryTar != "") {
		resultFiles, _ := utils.Untar(result)
		resultFiles.SaveAt(path)
	} else if(tryZip != "") {
		resultFiles, _ := utils.Unzip(result)
		resultFiles.SaveAt(path)
	} 

	return path, nil
}

func GetDependencyList(config map[string]interface {}) []map[string]interface {} {
	depList := (config["dependencies"].(map[string]interface {}))["default"].(map[string]interface {})
	result := make([]map[string]interface {}, 0)
	for name, value := range depList {
		dep := utils.BuildDependencyFromString("gupm", name)
		if(reflect.TypeOf(value).String() == "string") {
			dep["version"] = value
		} else {
			valueObject := value.(map[string]interface {})
			if(valueObject["provider"].(string) != "") {
				dep["provider"] = valueObject["provider"]
			}
			if(valueObject["version"].(string) != "") {
				dep["version"] = valueObject["version"]
			}
		}
		result = append(result, dep)
	}
	return result
}