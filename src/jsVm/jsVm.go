package jsVm

import (
	"../utils"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/Masterminds/semver"
)

func Setup(vm *otto.Otto) {	
	vm.Set("httpGet", func(call otto.FunctionCall) otto.Value {
		url, _ := call.Argument(0).ToString()
		res := utils.HttpGet(url)
		result, _ :=  vm.ToValue(utils.StringToJSON(string(res)))
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
				fmt.Println(err)
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
