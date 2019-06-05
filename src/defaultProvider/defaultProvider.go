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

func Bootstrap(path string) {
	if(utils.FileExists(utils.Path(path + "/gupm.json"))) {
		ui.Error("A project already exists in this folder. Aborting bootstrap.")
		return
	}

	name := ui.WaitForInput("Please enter the name of the project: ")
	description := ui.WaitForInput("Enter a description: ")
	author := ui.WaitForInput("Enter the author: ")
	licence := ui.WaitForInput("Enter the licence (ISC): ")
	
	if (name == "") {
		ui.Error("Name cannot be empty. Try again")
		return
	} else {
		if(licence == "") {
			licence = "ISC"
		}

		fileContent := `{
	"name": "` + name + `",
	"name": "0.0.1",
	"description": "` + description + `",
	"author": "` + author + `",
	"licence": "` + licence + `"
}`
		ioutil.WriteFile(utils.Path(path + "/gupm.json"), []byte(fileContent), os.ModePerm)
	}
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

	utils.SaveLockDep(path)

	return path, nil
}

func GetDependencyList(config map[string]interface {}) []map[string]interface {} {
	if(config == nil) {
		ui.Error("no config found. Please bootstrap the project with `g bootstrap`")
		return nil
	}
	depEnv, ok := config["dependencies"].(map[string]interface {})
	if(!ok) {
		ui.Log("no dependencies")
		return nil
	}
	depList := depEnv["default"].(map[string]interface {})
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

func ExpandDependency(dependency map[string]interface {}) (map[string]interface {}, error) {
	config := GetPackageConfig(utils.Path(dependency["path"].(string) + "/gupm.json"))
	dependency["dependencies"] = make(map[string]interface {})

	for _, depRaw := range (config["dependencies"].(map[string]interface {}))["default"].(map[string]interface {}) {
		dep := depRaw.(map[string]interface {})
		depBlock := utils.BuildDependencyFromString(dep["provider"].(string), dep["name"].(string))
		depBlock["version"] = dep["version"]
		dependency["dependencies"] = append(dependency["dependencies"].([]map[string]interface {}), depBlock)
	}

	return dependency, nil
}

func BinaryInstall(path string, packagePath string) error {
	packages := utils.ReadDir(packagePath)

	for _, dep := range packages {
		config := GetPackageConfig(utils.Path(packagePath + "/" + dep.Name() + "/gupm.json"))
		bins := config["binaries"].(map[string]string)
		for name, relPath := range bins {
			os.Symlink(utils.Path("../gupm_modules/" + "/" + dep.Name()  + relPath), ".bin/" + name)
		} 
	}

	return nil
}

func SaveDependencyList(depList []map[string]interface{}) error {
	config := GetPackageConfig("gupm.json")
	config["dependencies"] = make(map[string]interface{})
	config["dependencies"].(map[string]interface{})["default"] = make(map[string]interface {})

	for _, dep := range depList {
		config["dependencies"].(map[string]interface{})["default"].(map[string]interface{})[dep["provider"].(string) + "://" + dep["name"].(string)] = dep["version"].(string)
	}

	utils.WriteJsonFile("gupm.json", config)

	return nil
}