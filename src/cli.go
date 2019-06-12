package main

import (
	"regexp"
	"strings"
	"strconv"
    "fmt"
	"./ui"
	"./utils"
)

type json = utils.Json

type Arguments map[string] string
func (a *Arguments) AsJson() json {
	res := utils.Json{}
	for i, v := range *a {
		res[i] = v
	}
	return res
}
func (a *Arguments) Join() string {
	res := ""
	for i, v := range *a {
		if(v == "true" || v == "false" ) {
			res += "--" + strings.ToLower(i) + " "
		} else if ok, _ := regexp.MatchString(`^\$\d+`, i); ok {
			res += v + " "
		} else {
			res += "--" + strings.ToLower(i) + " " + v + " "
		}
	}
	return strings.TrimSpace(res)
}
func (a *Arguments) AsList() []string {
	res := []string{}
	for i, v := range *a {
		if ok, _ := regexp.MatchString(`^\$\d+`, i); ok && i != "$0" {
			res = append(res, v)
		}
	}
	return res
}

func GetArgs(args []string) (string, Arguments) {
	arguments := make(Arguments)
	next := ""
	dolsI := 1
	
	if(len(args) == 0) {
		arguments["$0"] = ""
		return "", arguments
	}

	command := args[0]
	if(len(args) < 2) {
		arguments["$0"] = command
		return command, arguments
	}

	argsToParse := args[1:]
	
	for _, value := range argsToParse {	
		nameCheck := regexp.MustCompile(`^--?(\w+)`)
		tryname := nameCheck.FindString(value)
		if(tryname != "") {
			long, _ := regexp.MatchString(`^--`, tryname)
			if(long) {
				tryname = tryname[1:]
			}

			if(next != "") {
				arguments[next] = "true"
				next = ""
			}
			next = strings.ToUpper(tryname[1:2]) + tryname[2:]			
		} else {
			if(next != "") {
				arguments[next] = value
				next = ""
			} else {
				arguments["$" + strconv.FormatInt(int64(dolsI), 10)] = value
				dolsI++
			}
		}
	}

	if(next != "") {
		arguments[next] = "true"
	}
	
	arguments["$0"] = command + " " + arguments.Join()

	return command, arguments
}

func getProvider(args Arguments) string {
	if args["Provider"] != "" {
		return args["Provider"]
	} else if(args["P"] != "") {
		return args["P"] 
	} else {
		return "gupm"
	}
}


func ExecCli(c string, args Arguments) (bool, error) {
	var err error
	notFound := "Cannot find commmand"
	
	if provider := getProvider(args); provider != ""  {
		Provider = provider
	}

	if c == "help" || c == "h" {
		fmt.Println("make / m :", "[--provider=]", "Install projects depdencies based on info in the entry point (depends on provider)")
		fmt.Println("install / i :", "[--provider=]", "Install package")
		fmt.Println("remove / r :", "[--provider=]", "remove package from module config")
		fmt.Println("publish / p :", "[--provider=]", "publish a project based on the model of your specific provider")
		fmt.Println("bootstrap / b :", "[--provider=]", "bootstrap a new project based on the model of your specific provider")

		fmt.Println("cache / c :", "clear or check the cache with \"cache clear\" or \"cache check\"")
		fmt.Println("self / s :", "self manage gupm. Try g \"self upgrade\" or \"g self uninstall\"")
		fmt.Println("plugin / pl :", "To install a plugin \"g pl install\". Then use \"g pl create\" to create a new one and \"g pl link\" to test your plugin")
	} else 

	if c == "make" || c == "m" {
		err = InstallProject(".")
	} else
	
	if c == "install" || c == "i" {
		err = AddDependency(".", args.AsList())
		if(err == nil) {
			err = InstallProject(".")
		}
	} else
	
	if c == "publish" || c == "p" {
		err = Publish(".")
	} else

	if c == "delete" || c == "d" {
		err = RemoveDependency(".", args.AsList())
	} else
	
	if c == "plugin" || c == "pl" {
		if(args["$1"] == "create") {
			PluginCreate(".")
		} else if (args["$1"] == "link") {
			PluginLink(".")
		} else if (args["$1"] == "install") {
			err = PluginInstall(".", args.AsList()[1:])
		} else if (args["$1"] == "delete") {
			PluginDelete(".", args.AsList()[1:])
		} else {
			ui.Error(notFound, args["$1"], "\n", "try cache clear or cache check")
		}
	} else
	
	if c == "cache" || c == "c" {
		if(args["$1"] == "clear") {
			CacheClear()
		} else if (args["$1"] == "check") {
			ui.Error("Not implemented yet.")
		} else {
			ui.Error(notFound, args["$1"], "\n", "try cache clear or cache check")
		}
	} else

	if c == "self" || c == "s" {
		if(args["$1"] == "upgrade") {
			SelfUpgrade()
		} else if (args["$1"] == "uninstall") {
			SelfUninstall()
		}else {
			ui.Error(notFound, args["$1"])
		}
	} else
	
	if c == "bootstrap" || c == "b" {
		err = Bootstrap(".")
	} else 
	
	if c == "test" || c == "t" {
		err = RunTest("tests")
	} else 
	
	{
		return false, nil
	}

	return true, err
}