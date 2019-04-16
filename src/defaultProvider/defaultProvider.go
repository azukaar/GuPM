package defaultProvider

import (
	"io/ioutil"
	"encoding/json"
	"../utils"
	"os"
	"regexp"
	"fmt"
)

func GetPackageConfig(entryPoint string) map[string]interface {} {
	var packageConfig map[string]interface{}
	b, err := ioutil.ReadFile(entryPoint)
	if err != nil {
		fmt.Println(err, entryPoint)
	}

	json.Unmarshal([]byte(string(b)), &packageConfig)

	return packageConfig
}

func GetDepedency(provider string, name string, version string, url string, path string) (string, error) {
	return string(utils.HttpGet(url)), nil
}

func PostGetDepedency(provider string, name string, version string, url string, path string, result string) (string, error) {
	os.MkdirAll(path, os.ModePerm)
	// resultBinary := []byte(result)
	extensionCheck := regexp.MustCompile(`\.tgz$`)
	tryExtension := extensionCheck.FindString(url)

	if(tryExtension != "") {
		resultFiles, _ := utils.Untar(result)
		resultFiles.SaveAt(path)
		// fmt.Println(resultFiles)
	} else {

	}

	return path, nil
}

func GetDependencyList(config map[string]interface {}) utils.PackageDepedencyListType {
	return config["dependencies"].(utils.PackageDepedencyListType)
}