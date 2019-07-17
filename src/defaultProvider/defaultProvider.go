package defaultProvider

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"

	"../ui"
	"../utils"
)

func Bootstrap(path string) {
	if utils.FileExists(utils.Path(path + "/gupm.json")) {
		ui.Error("A project already exists in this folder. Aborting bootstrap.")
		return
	}

	name := ui.WaitForInput("Please enter the name of the project: ")
	description := ui.WaitForInput("Enter a description: ")
	author := ui.WaitForInput("Enter the author: ")
	licence := ui.WaitForInput("Enter the licence (ISC): ")

	if name == "" {
		ui.Error("Name cannot be empty. Try again")
		return
	} else {
		if licence == "" {
			licence = "ISC"
		}

		fileContent := `{
	"name": "` + name + `",
	"version": "0.0.1",
	"description": "` + description + `",
	"author": "` + author + `",
	"licence": "` + licence + `"
}`
		ioutil.WriteFile(utils.Path(path+"/gupm.json"), []byte(fileContent), os.ModePerm)
	}
}

func GetPackageConfig(entryPoint string) map[string]interface{} {
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
	gzCheck := regexp.MustCompile(`\.gz$`)
	trygz := gzCheck.FindString(url)
	zipCheck := regexp.MustCompile(`\.zip$`)
	tryZip := zipCheck.FindString(url)

	if tryTar != "" {
		resultFiles, err := utils.Untar(result)
		if err != nil {
			return path, err
		}
		resultFiles.SaveAt(path)
	} else if trygz != "" {
		resultFiles, err := utils.Ungz(result)
		if err != nil {
			return path, err
		}
		resultFiles.SaveAt(path)
	} else if tryZip != "" {
		resultFiles, err := utils.Unzip(result)
		if err != nil {
			return path, err
		}
		resultFiles.SaveAt(path)
	}

	utils.SaveLockDep(path)

	return path, nil
}

func GetDependencyList(config map[string]interface{}) []map[string]interface{} {
	if config == nil {
		ui.Error("no config found. Please bootstrap the project with `g bootstrap`")
		return nil
	}
	depEnv, ok := config["dependencies"].(map[string]interface{})
	if !ok {
		ui.Log("no dependencies")
		return nil
	}
	depList, hasDefault := depEnv["default"].(map[string]interface{})
	if !hasDefault {
		ui.Log("no dependencies")
		return nil
	}
	result := make([]map[string]interface{}, 0)
	for name, value := range depList {
		dep := utils.BuildDependencyFromString("gupm", name)
		if reflect.TypeOf(value).String() == "string" {
			dep["version"] = value
		} else {
			valueObject := value.(map[string]interface{})
			if valueObject["provider"].(string) != "" {
				dep["provider"] = valueObject["provider"]
			}
			if valueObject["version"].(string) != "" {
				dep["version"] = valueObject["version"]
			}
		}
		result = append(result, dep)
	}
	return result
}

func ExpandDependency(dependency map[string]interface{}) (map[string]interface{}, error) {
	configFilePath := utils.Path(dependency["path"].(string) + "/gupm.json")

	if utils.FileExists(configFilePath) {
		config := GetPackageConfig(configFilePath)
		dependency["dependencies"] = make(map[string]interface{})

		if config["dependencies"] != nil {
			dependency["dependencies"] = GetDependencyList(config)
		}
	}

	return dependency, nil
}

func BinaryInstall(path string, packagePath string) error {
	packages, _ := utils.ReadDir(packagePath)

	for _, dep := range packages {
		configFilePath := utils.Path(packagePath + "/" + dep.Name() + "/gupm.json")
		if utils.FileExists(configFilePath) {
			config := GetPackageConfig(configFilePath)
			if config["binaries"] != nil {
				bins := config["binaries"].(map[string]string)
				for name, relPath := range bins {
					os.Symlink(utils.Path("../gupm_modules/"+"/"+dep.Name()+relPath), ".bin/"+name)
				}
			}
		}
	}

	return nil
}

func SaveDependencyList(depList []map[string]interface{}) error {
	config := GetPackageConfig("gupm.json")
	if config["dependencies"] == nil {
		config["dependencies"] = make(map[string]interface{})
	}
	config["dependencies"].(map[string]interface{})["default"] = make(map[string]interface{})

	for _, dep := range depList {
		key := utils.BuildStringFromDependency(map[string]interface{}{
			"provider": dep["provider"].(string),
			"name":     dep["name"].(string),
		})

		config["dependencies"].(map[string]interface{})["default"].(map[string]interface{})[key] = dep["version"].(string)
	}

	utils.WriteJsonFile("gupm.json", config)

	return nil
}
