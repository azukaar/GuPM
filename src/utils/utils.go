package utils

import (
	"encoding/json"
	"os"
)

type Json map[string]interface {}
type JsonMap []interface {}
type PackageDepedencyListType []map[string]interface {}

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