package utils

import (
	"encoding/json"
	"os"
	"net/http"
	"regexp"
	"fmt"
	"io/ioutil"
	"github.com/otiai10/copy"
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

func JsonExport(input otto.Value) map[string] interface {} {
	var obj map[string] interface {}
	result := make(map[string] interface {})
	exported, _ := input.Export()
	obj, _ = exported.(map[string]interface {})
	for index, value := range obj {
		ottoValue, _ := otto.ToValue(value)
		if (otto.Value.IsObject(ottoValue)) {
			result[index] = JsonExport(ottoValue)
		} else {
			result[index] = value
		}
	}

	return result
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

func CopyRecursive(dest string, source string) {
	copy.Copy(source, dest)
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