package main

import "encoding/json"
import "net/http"

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"
)

func download(packagename string, version string, ch chan<-string) {
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
		fmt.Println("for", packagename)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("!!!", err)
		fmt.Println("for", packagename)
	}

	os.MkdirAll("node_modules/"+packagename,  os.ModePerm)

	err = Untar("node_modules/"+packagename, string(body))
	if err != nil {
		fmt.Println(err)
		fmt.Println("node_modules/"+packagename, version)
	}
	
	ch <- url
}

func depInstall(workingDir string, file string) {
	ch := make(chan string)
	var packagejson map[string]interface{}
	b, err := ioutil.ReadFile(workingDir+"package.json") // just pass the file name
	if err != nil {
			fmt.Println(err)
	}

	json.Unmarshal([]byte(string(b)), &packagejson)

	if(packagejson["dependencies"] != nil) {
		var dependencies = packagejson["dependencies"].(map[string]interface{})
		var newDeps = make(map[string]interface{})

		fmt.Println("-- installing dependencies for", workingDir, file);

		for packagename, versionBlob := range dependencies {
			if _, err := os.Stat("node_modules/"+packagename); err != nil {
				var version string = versionBlob.(string)
				if(version[0:1] == "^" || version[0:1] == "~") {
					version = versionBlob.(string)[1:]
				}
				go download(packagename, version, ch)
				newDeps[packagename] = version
			}
		}

		for range newDeps {
			<-ch
		}

		for packagename, version := range newDeps {
			_ = version
			depInstall("node_modules/"+packagename+"/", "package.json")
		}
	}
}

func main() {
	start := time.Now()
	depInstall("", "package.json")
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}