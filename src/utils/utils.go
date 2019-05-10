package utils

import (
	"encoding/json"
	"os"
	"net/http"
	"regexp"
	"fmt"
	"io/ioutil"
	"github.com/robertkrimen/otto"
    "path/filepath"
)

type Json map[string]interface {}
type JsonMap []interface {}
type PackageDependencyListType []map[string]interface {}

type Dependency struct {
	Name string
	Provider string
	Version string
}

func JsonExport(input interface {}) interface {} {
	asMap, isMap := input.(map[string] interface {})
	asSlice, isSlice := input.([] interface {})
	if(isMap) {
		for index, value := range asMap {
			asValue, ok := value.(otto.Value)
			if(ok) {
				exported, _ := asValue.Export()
				asMap[index] = JsonExport(exported)
			}
		}
		return asMap
	} else if(isSlice) {
		for index, value := range asSlice {
			asValue, ok := value.(otto.Value)
			if(ok) {
				exported, _ := asValue.Export()
				asSlice[index] = JsonExport(exported)
			}
		}
		return asSlice
	} else {
		return input
	}
}

func BuildDependencyFromString(defaultProvider string, dep string) map[string]interface {} {
	result := make(map[string]interface {})
	step := dep

	versionCheck := regexp.MustCompile(`@[\w\.\-\_\^\~]+$`)
	tryversion := versionCheck.FindString(step)
	if(tryversion != "") {
		result["version"] = tryversion[1:]
		step = versionCheck.ReplaceAllString(step, "")
	} else {
		result["version"] = "*.*.*"
	}

	providerCheck := regexp.MustCompile(`^[\w\-\_]+\:\/\/`)
	tryprovider := providerCheck.FindString(step)
	if(tryprovider != "") {
		result["provider"] = tryprovider[:len(tryprovider)-3]
		step = providerCheck.ReplaceAllString(step, "")
	} else {
		result["provider"] = defaultProvider
	}

	result["name"] = step
	return result
}

func StringToJSON(b string) map[string]interface {} {
	var jsonString map[string]interface{}
	json.Unmarshal([]byte(string(b)), &jsonString)
	return jsonString
}

func ReadJSON(path string, target interface{}) error  {
	b, err := os.Open(path) // just pass the file name
	if err != nil {
		return err
	}

	return json.NewDecoder(b).Decode(target)
}

func HttpGet(url string) []byte {
	resp, httperr := http.Get(url)
	if httperr != nil {
		fmt.Println("Error accessing ", url, " trying again. Check your network.")
		return HttpGet(url)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading HTTP response", err)
		return HttpGet(url)
	}
	
	return body
}

func FileExists(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return true
}

func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func RemoveIndex(s []map[string]interface {}, index int) []map[string]interface {} {
    return append(s[:index], s[index+1:]...)
}

func ReadDir(path string) []os.FileInfo{
    files, err := ioutil.ReadDir(path)
    if err != nil {
        fmt.Println(err)
	}

    return files
}

func DIRNAME() string {
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    dir := filepath.Dir(ex)
	return dir
}