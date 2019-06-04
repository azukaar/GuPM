package jsVm

import (
	"../utils"
	"../ui"
	"errors"
	"sync"
	"io/ioutil"
	"os"
	"runtime"
	"time"
	"encoding/json"
	"github.com/robertkrimen/otto"
	"github.com/Masterminds/semver"
)

var lock = sync.RWMutex{}
var scriptCache = make(map[string]string)

func Run(path string, input map[string]interface {}) (otto.Value, error) {
	var err error
	var ret otto.Value

	lock.Lock()
	if(scriptCache[path] == "") {
		file, err := ioutil.ReadFile(path)
		if(err != nil) {
			return otto.UndefinedValue(),  err
		}
		scriptCache[path] = string(file)
	}
	script := scriptCache[path]
	lock.Unlock()

	vm := otto.New()
	vm.Interrupt = make(chan func(), 1) 
	Setup(vm)

	for varName, varValue := range input {
		vm.Set(varName, varValue)
	}

	ret, err = vm.Run(script)

	if(err != nil) {
		ui.Error(err.Error())
		return otto.UndefinedValue(),  errors.New("Error occured while executing the GS code")
	}

	return ret, nil
}

func Setup(vm *otto.Otto) {	
	vm.Set("httpGetJson", func(call otto.FunctionCall) otto.Value {
		url, _ := call.Argument(0).ToString()
		res := utils.HttpGet(url)
		result, _ :=  vm.ToValue(utils.StringToJSON(string(res)))
		return result
	})

	vm.Set("httpGet", func(call otto.FunctionCall) otto.Value {
		url, _ := call.Argument(0).ToString()
		res := utils.HttpGet(url)
		result, _ :=  vm.ToValue(string(res))
		return result
	})

	vm.Set("dir", func(call otto.FunctionCall) otto.Value {
		glob, _ := call.Argument(0).ToString()
		res, _ := utils.Dir(glob)
		result, _ :=  vm.ToValue(res)
		return result
	})
	
	vm.Set("osSleep", func(call otto.FunctionCall) otto.Value {
		timeMs, _ := call.Argument(0).ToInteger()
		time.Sleep(time.Duration(timeMs) * time.Millisecond)
		result, _ :=  vm.ToValue(true)
		return result
	})

	vm.Set("readJsonFile", func(call otto.FunctionCall) otto.Value {
		path, _ := call.Argument(0).ToString()
		path = utils.Path(path)
		b, err := ioutil.ReadFile(path)
		if(err != nil) {
			ui.Error(err.Error())
		}
		result, _ :=  vm.ToValue(utils.StringToJSON(string(b)))
		return result
	})

	vm.Set("readFile", func(call otto.FunctionCall) otto.Value {
		path, _ := call.Argument(0).ToString()
		path = utils.Path(path)
		b, err := ioutil.ReadFile(path)
		if(err != nil) {
			ui.Error(err.Error())
		}
		result, _ :=  vm.ToValue(string(b))
		return result
	})

	vm.Set("removeFiles", func(call otto.FunctionCall) otto.Value {
		files, _ := call.Argument(0).Export()
		_, isString := files.(string)
		if(isString) {
			files = []string{files.(string)}
		}
		utils.RemoveFiles(files.([]string))
		result, _ :=  vm.ToValue(true)
		return result
	})

	vm.Set("copyFiles", func(call otto.FunctionCall) otto.Value {
		files, _ := call.Argument(0).Export()
		_, isString := files.(string)
		if(isString) {
			files = []string{files.(string)}
		}
		path, _ := call.Argument(1).ToString()
		path = utils.Path(path)
		utils.CopyFiles(files.([]string), path)
		result, _ :=  vm.ToValue(true)
		return result
	})

	vm.Set("pwd", func(call otto.FunctionCall) otto.Value {
		dir, _ := os.Getwd()
		result, _ :=  vm.ToValue(dir)
		return result
	})

	vm.Set("env", func(call otto.FunctionCall) otto.Value {
		name, _ := call.Argument(0).ToString()
		value, _ := call.Argument(1).ToString()

		if(value == "undefined") {
			result, _ :=  vm.ToValue(os.Getenv(name))
			return result
		} else {
			os.Setenv(name, value)
			res, _ := vm.ToValue(true)
			return res
		}
	})

	vm.Set("exec", func(call otto.FunctionCall) otto.Value {
		exec, _ := call.Argument(0).ToString()
		args, _ := call.Argument(1).Export()
		_, ok := args.([]string)
		if(!ok) {
			args = make([]string, 0)
		}

		utils.ExecCommand(exec, args.([]string))

		result, _ :=  vm.ToValue(true)
		return result
	})

	vm.Set("run", func(call otto.FunctionCall) otto.Value {
		exec, _ := call.Argument(0).ToString()
		args, _ := call.Argument(1).Export()
		_, ok := args.([]string)
		if(!ok) {
			args = make([]string, 0)
		}

		res, _ := utils.RunCommand(exec, args.([]string))
		res = res[:len(res)-1]

		result, _ :=  vm.ToValue(res)
		return result
	})
	
	vm.Set("exit", func() {
		os.Exit(1)
	})

	vm.Set("writeJsonFile", func(call otto.FunctionCall) otto.Value {
		path, _ := call.Argument(0).ToString()
		path = utils.Path(path)
		toExport, _ := call.Argument(1).Export()
		file := JsonExport(toExport).(map[string] interface {})
		bytes, _ := json.MarshalIndent(file, "", "    ")
		err := ioutil.WriteFile(path, bytes, os.ModePerm)
		if(err != nil) {
			ui.Error(err.Error())
		}
		result, _ :=  vm.ToValue(true)
		return result
	})

	vm.Set("writeFile", func(call otto.FunctionCall) otto.Value {
		path, _ := call.Argument(0).ToString()
		path = utils.Path(path)
		toExport, _ := call.Argument(1).ToString()
		err := ioutil.WriteFile(path, []byte(toExport), os.ModePerm)
		if(err != nil) {
			ui.Error(err.Error())
		}
		result, _ :=  vm.ToValue(true)
		return result
	})

	vm.Set("_OSNAME", runtime.GOOS)

	vm.Set("mkdir", func(call otto.FunctionCall) otto.Value {
		path, _ := call.Argument(0).ToString()
		path = utils.Path(path)
		os.MkdirAll(path, os.ModePerm)
		result, _ :=  vm.ToValue(true)
		return result
	})
	
	vm.Set("saveLockDep", func(call otto.FunctionCall) otto.Value {
		path, _ := call.Argument(0).ToString()
		path = utils.Path(path)
		utils.SaveLockDep(path)
		result, _ :=  vm.ToValue(true)
		return result
	})
	
	vm.Set("fileExists", func(call otto.FunctionCall) otto.Value {
		path, _ := call.Argument(0).ToString()
		path = utils.Path(path)
		res := utils.FileExists(path)
		result, _ :=  vm.ToValue(res)
		return result
	})

	vm.Set("waitForInput", func(call otto.FunctionCall) otto.Value {
		msg, _ := call.Argument(0).ToString()
		res := ui.WaitForInput(msg)
		result, _ :=  vm.ToValue(res)
		return result
	})

	vm.Set("waitForMenu", func(call otto.FunctionCall) otto.Value {
		msg, _ := call.Argument(0).Export()
		res := ui.WaitForMenu(msg.([]string))
		result, _ :=  vm.ToValue(res)
		return result
	})

	vm.Set("waitForKey", func(call otto.FunctionCall) otto.Value {
		ui.WaitForKey()
		result, _ :=  vm.ToValue(true)
		return result
	})

	vm.Set("tar", func(call otto.FunctionCall) otto.Value {
		files, _ := call.Argument(0).Export()
		_, isString := files.(string)
		if(isString) {
			files = []string{files.(string)}
		}
		res, err := utils.Tar(files.([]string))
		if(err != nil) {
			ui.Error(err.Error())
		}
		b, _ := json.Marshal(res)
		result, _ :=  vm.ToValue(utils.StringToJSON(string(b)))
		return result
	})

	vm.Set("readDir", func(call otto.FunctionCall) otto.Value {
		path, _ := call.Argument(0).ToString()
		path = utils.Path(path)
		var filenames = make([]string, 0)
		files := utils.ReadDir(path)
		for _, file := range files {
			filenames = append(filenames, file.Name())
		}
		result, _ :=  vm.ToValue(filenames)
		return result
	})

	vm.Set("createSymLink", func(call otto.FunctionCall) otto.Value {
		from, _ := call.Argument(0).ToString()
		from = utils.Path(from)
		to, _ := call.Argument(1).ToString()
		to = utils.Path(to)
		err := os.Symlink(from, to)
		if(err != nil) {
			ui.Error(err.Error())
		}
		result, _ :=  vm.ToValue(true)
		return result
	})

	vm.Set("untar", func(call otto.FunctionCall) otto.Value {
		var res utils.FileStructure
		file, _ := call.Argument(0).ToString()
		res, _ = utils.Untar(file)
		b, _ := json.Marshal(res)
		result, _ :=  vm.ToValue(utils.StringToJSON(string(b)))
		return result
	})

	vm.Set("unzip", func(call otto.FunctionCall) otto.Value {
		var res utils.FileStructure
		file, _ := call.Argument(0).ToString()
		res, _ = utils.Unzip(file)
		b, _ := json.Marshal(res)
		result, _ :=  vm.ToValue(utils.StringToJSON(string(b)))
		return result
	})

	vm.Set("saveFileAt", func(call otto.FunctionCall) otto.Value {
		var fs utils.FileStructure
		file, _ := call.Argument(0).Export()
		path, _ := call.Argument(1).ToString()
		path = utils.Path(path)
		bytes, _ := json.Marshal(file)
		json.Unmarshal(bytes, &fs)
		fs.SaveAt(path)
		result, _ := vm.ToValue(path)
		return result
	})

	vm.Set("semverInRange", func(call otto.FunctionCall) otto.Value {
		rangeStr, _ := call.Argument(0).ToString()
		version, _ := call.Argument(1).ToString()
		rangeVer, _ := semver.NewConstraint(rangeStr)
		sver, _ := semver.NewVersion(version)
		value := rangeVer.Check(sver)
		result, _ :=  vm.ToValue(value)
		return result
	})

	vm.Set("semverLatestInRange", func(call otto.FunctionCall) otto.Value {
		rangeStr, _ := call.Argument(0).ToString()
		versionList, _ := call.Argument(1).Export()
		var version string
		var versionSem *semver.Version
		rangeVer, _ := semver.NewConstraint(rangeStr)

		for _, verCand := range versionList.([]string) {
			sver, err := semver.NewVersion(verCand)
			if err != nil {
				ui.Error(err.Error())
			}
			
			if(rangeVer.Check(sver) && (versionSem == nil || sver.GreaterThan(versionSem))) {
				version = verCand
				versionSem = sver
			}
		}
		if(version != "") {
			result, _ :=  vm.ToValue(version)
			return result
		} else {
			return otto.UndefinedValue()
		}
	})
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