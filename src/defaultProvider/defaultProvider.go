package defaultProvider

import (
	"io/ioutil"
	"encoding/json"
	"../utils"
	"os"
	"regexp"
	"fmt"
)

func SaveDependencyList(depList []map[string]interface{}) error {
	fmt.Println(depList)
	return nil
}

func GetPackageConfig(entryPoint string) map[string]interface {} {
	var packageConfig map[string]interface{}
	b, err := ioutil.ReadFile(entryPoint)
	if err != nil {
		fmt.Println(err, entryPoint)
	}

	json.Unmarshal([]byte(string(b)), &packageConfig)

	return packageConfig
}

func GetDependency(provider string, name string, version string, url string, path string) (string, error) {
	return string(utils.HttpGet(url)), nil
}

func PostGetDependency(provider string, name string, version string, url string, path string, result string) (string, error) {
	os.MkdirAll(path, os.ModePerm)
	extensionCheck := regexp.MustCompile(`\.tgz$`)
	tryExtension := extensionCheck.FindString(url)

	if(tryExtension != "") {
		resultFiles, _ := utils.Untar(result)
		resultFiles.SaveAt(path)
	} else {

	}

	return path, nil
}

func GetDependencyList(config map[string]interface {}) []map[string]interface {} {
	// depList := (config["dependencies"].(map[string]interface {}))["default"].(map[string]interface {})
	result := make([]map[string]interface {}, 0)
	// for name, version := range depList {
	// 	result.append()
	// }
	return result
}