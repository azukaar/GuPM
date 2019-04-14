package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func readJSON(path string) map[string]interface{} {
	var jsonfile map[string]interface{}
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal([]byte(string(b)), &jsonfile)

	return jsonfile
}