package utils

import (
	"encoding/json"
	"os"
	"net/http"
	"fmt"
	"io/ioutil"
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

func FileExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

