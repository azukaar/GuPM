package main

import "encoding/json"
import "net/http"

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"regexp"
	"github.com/blang/semver"
	"strings"
)

func httpGet(url string) []byte {
	resp, httperr := http.Get(url)
	if httperr != nil {
		fmt.Println("Error trying to dl file ", url, " trying again. Check your network.")
		return httpGet(url)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("!", err)
	}
	
	return body
}

func is404(body []byte) bool {
	return string(body) == "{\"error\":\"Not found\"}"
}

func resolveTag(packagename string, tag string) string {
	var url = "https://registry.npmjs.org/" + packagename
	var result map[string]interface{}
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
	
	json.Unmarshal([]byte(body), &result)
	var tags = result["dist-tags"].(map[string]interface{})

	if(tags != nil && tags[tag] != nil) {
		return (tags[tag].(string))
	}
	
	return tag
}

func download(packagename string, versionBlob string, ch chan<-string) {
	var version, correctVersion, url string
	var body []byte
	var packageDescription map[string]interface{}
	
	specificCheck := regexp.MustCompile(`^\d+\.\d+\.\d+\-[\w\d\.]+$`)
	trySpecific := specificCheck.FindString(versionBlob)

	if(trySpecific != "") {
		url = "https://registry.npmjs.org/" + 
				packagename +
				"/-/" +
				packagename +
				"-" +
				versionBlob +
				".tgz"
		body = httpGet(url)
	} else {
		res := httpGet("https://registry.npmjs.org/" + packagename)
		json.Unmarshal([]byte(res), &packageDescription)
		candidates := packageDescription["versions"].(map[string]interface{})
	
		correctVersion = strings.Replace(versionBlob, "~", ">=", -1)
		correctVersion = strings.Replace(correctVersion, "^", ">=", -1)
		correctVersion = strings.Replace(correctVersion, " -", " <", -1)
		correctVersion = strings.Replace(correctVersion, "*", ">0.0.0", -1)
	
		correctVersion = regexp.MustCompile(`^([\>\<\=]{0,2}\d{1,3}\.\d{1,3})$`).ReplaceAllString(correctVersion, "$1.0")
		correctVersion = regexp.MustCompile(`^([\>\<\=]{0,2}\d{1,3})$`).ReplaceAllString(correctVersion, "$1.0.0")
	
		correctVersion = regexp.MustCompile(`^([\>\<\=]{0,2}\d{1,3}\.\d{1,3})\s`).ReplaceAllString(correctVersion, "$1.0")
		correctVersion = regexp.MustCompile(`^([\>\<\=]{0,2}\d{1,3})\s`).ReplaceAllString(correctVersion, "$1.0.0")
	
		correctVersion = regexp.MustCompile(`[\s\>\<\=]([\>\<\=]{0,2}\d{1,3}\.\d{1,3})$`).ReplaceAllString(correctVersion, "$1.0")
		correctVersion = regexp.MustCompile(`[\s\>\<\=]([\>\<\=]{0,2}\d{1,3})$`).ReplaceAllString(correctVersion, "$1.0.0")
	
		correctVersion = regexp.MustCompile(`^(\d+\.\d+)\.x$`).ReplaceAllString(correctVersion, ">$1.0")
		correctVersion = regexp.MustCompile(`^(\d+)\.x\.x$`).ReplaceAllString(correctVersion, ">$1.0.0")
	
		rangeVer, err := semver.ParseRange(correctVersion)
		if err != nil {
			fmt.Println("ERR 2", err)
			fmt.Println("ERR 2", correctVersion + "(" + versionBlob + ")")
			fmt.Println("ERR 2", packagename)
		}

		url = "NO MATCHING VERSION FOR " + packagename + " " +  correctVersion
	
		for verCand := range candidates {
			sver, err := semver.ParseTolerant(verCand)
			if err != nil {
				fmt.Println("!", err)
			}
	
			if(rangeVer(sver)) {
				version = verCand
				url = "https://registry.npmjs.org/" + 
						packagename +
						"/-/" +
						packagename +
						"-" +
						version +
						".tgz"
				body = httpGet(url)
				break;
			}
		}
	}

	if (is404(body)) {
		fmt.Println("404 NOT FOUND, " + url)
	} else {
		fmt.Println(url)
		os.MkdirAll("node_modules/"+packagename,  os.ModePerm)
	
		err := Untar("node_modules/"+packagename, string(body))
		if err != nil {
			fmt.Println("1", err)
			fmt.Println("node_modules/"+packagename, version)
		}	
	}

	ch <- packagename
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

		fmt.Println("-- installing dependencies for", workingDir);

		for packagename, versionBlob := range dependencies {
			if _, err := os.Stat("node_modules/"+packagename); err != nil {
				var version string = versionBlob.(string)
				tagCheck := regexp.MustCompile(`^\d*_*\w+[\d\w_]*$`)
				tryTag := tagCheck.FindString(version)

				if (tryTag != "") {
					version = resolveTag(packagename, tryTag)
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