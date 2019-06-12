package utils

import (
	"encoding/json"
	"os"
	"net/http"
	"regexp"
	"time"
	"../ui"
	"os/exec"
	"runtime"
	"reflect"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
    "path/filepath"
	"gopkg.in/go-playground/validator.v9"
)

type Dependency struct {
	Name string
	Provider string
	Version string
}

func buildCmd(toRun string, args []string) *exec.Cmd{
	isNode := regexp.MustCompile(`.js$`)
	var cmd *exec.Cmd
	bashargs := []string{}

	// temporary hack to make windows execute js file with node
	if(isNode.FindString(toRun) != "") {
		bashargs = append(bashargs, toRun)
		bashargs = append(bashargs, args...)
		cmd = exec.Command("node", bashargs...)	
	} else {
		bashargs = append(bashargs, args...)
		cmd = exec.Command(toRun, bashargs...)	
	}

	return cmd
}

func ExecCommand(toRun string, args []string) error {
	cmd := buildCmd(toRun, args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
	return nil
}

func ReadGupmJson(path string) (*GupmEntryPoint, error) {
	validate := validator.New()

	config := new(GupmEntryPoint)
	errRead := ReadJSON(path, config)
	if(errRead != nil) {
		ui.Error("Could not find", path)
		return nil, errRead
	}
	errValidate := validate.Struct(config)
	if(errValidate != nil) {
		ui.Error("Error validating ", path)
		return nil, errValidate
	}
	return config, nil
}

func RunCommand(toRun string, args []string) (string, error) {
	cmd := buildCmd(toRun, args)
	res, err := cmd.Output()
	if(err != nil) {
		return "", err
	}
	return string(res), nil
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

var numberConnectionOpened = 0

func HttpGet(url string) []byte {
	if(numberConnectionOpened > 50) {
		time.Sleep(1000 * time.Millisecond)
		return HttpGet(url)
	}

	numberConnectionOpened++
	resp, httperr := http.Get(url)
	if httperr != nil {
		numberConnectionOpened--
		isRateLimit, _ := regexp.MatchString(`unexpected EOF$`, httperr.Error())
		if(!isRateLimit) {
			ui.Error("Error accessing", url, "trying again.", httperr)
		}
		time.Sleep(1000 * time.Millisecond)
		return HttpGet(url)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		numberConnectionOpened--
		ui.Error("Error reading HTTP response ", err)
		return HttpGet(url)
	}
	
	numberConnectionOpened--
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

func RecursiveFileWalkDir(source string) []string {
	result := make([]string, 0)
	err := filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if(!info.IsDir()){
				result = append(result, path)
			}
			return nil
		})
	if err != nil {
		ui.Error(err)
	}
	return result
}

func ReadDir(path string) []os.FileInfo{
    files, err := ioutil.ReadDir(path)
    if err != nil {
        ui.Error(err)
	}

    return files
}

func IsDirectory(path string) (bool) {
    fileInfo, err := os.Stat(path)
    if err != nil {
      return false
    }
    return fileInfo.IsDir()
}

func HOMEDIR(fallback string) string {
	hdir, errH := homedir.Dir()
	if(errH != nil) {
		ui.Error(errH)
		hdir = fallback
	}
	return hdir
}

func DIRNAME() string {
	ex, err := os.Executable()
	exr, err := filepath.EvalSymlinks(ex)
    if err != nil {
        panic(err)
    }
    dir := filepath.Dir(exr)
	return dir
}

func WriteFile(path string, file string) error {
	return ioutil.WriteFile(Path(path), []byte(file), os.ModePerm)
}

func WriteJsonFile(path string, file map[string]interface {}) {
	bytes, _ := json.MarshalIndent(file, "", "    ")
	err := ioutil.WriteFile(path, bytes, os.ModePerm)
	if(err != nil) {
		ui.Error(err)
	}
}

// TODO: https://blog.golang.org/pipelines
// add proper checksum check 

func SaveLockDep(path string) {
	ioutil.WriteFile(path+"/.gupm_locked", []byte("1"), os.ModePerm)
}

func AbsPath(rel string) string {
	abs, _ := filepath.Abs(rel)
	return abs
}

func Path(path string) string {
	if runtime.GOOS == "windows" {
		return filepath.FromSlash(path)
	} else {
		return filepath.ToSlash(path)
	}
}

func Contains(s interface{}, elem interface{}) bool {
    arrV := reflect.ValueOf(s)

    if arrV.Kind() == reflect.Slice {
        for i := 0; i < arrV.Len(); i++ {

            // XXX - panics if slice element points to an unexported struct field
            // see https://golang.org/pkg/reflect/#Value.Interface
            if arrV.Index(i).Interface() == elem {
                return true
            }
        }
    }

    return false
}


func ArrString(something interface{}) []string {
	_, ok := something.([]string) 
	if(!ok) {
		res := make([]string, 0)
		for _, v := range something.([]interface{}) {
			res = append(res, v.(string))
		}
		return res
	} else {
		return something.([]string)
	}
}