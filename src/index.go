package main

import "encoding/json"
import "net/http"

import (
	"os"
	"strings"
	"fmt"
	"time"
	"io/ioutil"
)

func download(root string, packagename string, version string, ch chan<-string) {
	var url = "https://registry.npmjs.org/" + 
						packagename +
						"/-/" +
						packagename +
						"-" +
						version +
						".tgz"

	resp, httperr := http.Get(url)
	if httperr != nil {
		fmt.Println("!!!", httperr)
		fmt.Println("for", root, packagename)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("!!!", err)
		fmt.Println("for", root, packagename)
	}

	var filename = packagename+".tgz"
	var folders = strings.Split(filename, "/")
	var folder = folders[:len(folders)-1]
	os.MkdirAll("temp/"+strings.Join(folder[:], "/"),  os.ModePerm)
	os.MkdirAll(root+"node_modules/"+packagename,  os.ModePerm)

	fileerr := ioutil.WriteFile("temp/"+filename, body, 0644)
	if fileerr != nil {
		fmt.Println("!!!", fileerr)
		fmt.Println("for", root, packagename)
	}

	err = Untar(root+"node_modules/"+packagename, string(body))
	if err != nil {
			fmt.Print(err)
	}
	
	ch <- url
}

func depInstall(workingDir string, file string) {
	fmt.Println("-- installing dependencies for", workingDir, file);
	ch := make(chan string)
	var packagejson map[string]interface{}
	b, err := ioutil.ReadFile(workingDir+"package.json") // just pass the file name
	if err != nil {
			fmt.Print(err)
	}

	json.Unmarshal([]byte(string(b)), &packagejson)

	if(packagejson["dependencies"] != nil) {
		var dependencies = packagejson["dependencies"].(map[string]interface{})

		for packagename, versionBlob := range dependencies {
			var version string = versionBlob.(string)
			if(version[0:1] == "^") {
				version = versionBlob.(string)[1:]
			}
			go download(workingDir, packagename, version, ch)
		}
		
		for range dependencies {
			fmt.Println(<-ch)
		}

		for packagename, version := range dependencies {
			_ = version
			depInstall(workingDir+"node_modules/"+packagename+"/", "package.json")
		}
	}
}

func main() {
	start := time.Now()
	depInstall("", "package.json")
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}