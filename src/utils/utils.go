package utils

import (
	"encoding/json"
	"os"
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/otiai10/copy"
)

type Json map[string]interface {}
type JsonMap []interface {}
type PackageDependencyListType []map[string]interface {}

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
		fmt.Println("Error trying to dl file ", url, " trying again. Check your network.")
		return HttpGet(url)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
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
