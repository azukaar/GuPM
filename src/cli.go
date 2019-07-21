package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"./ui"
	"./utils"
)

type json = utils.Json

var ProviderWasForced = false

type Arguments map[string]string

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
		if v == "true" || v == "false" {
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
	i := 1

	for (*a)["$"+strconv.Itoa(i)] != "" {
		res = append(res, (*a)["$"+strconv.Itoa(i)])
		i++
	}

	return res
}

func GetArgs(args []string) (string, Arguments) {
	arguments := make(Arguments)
	next := ""
	dolsI := 1

	if len(args) == 0 {
		arguments["$0"] = ""
		return "", arguments
	}

	command := args[0]
	if len(args) < 2 {
		arguments["$0"] = command
		return command, arguments
	}

	argsToParse := args[1:]

	for _, value := range argsToParse {
		nameCheck := regexp.MustCompile(`^--?(\w+)`)
		tryname := nameCheck.FindString(value)
		if tryname != "" {
			long, _ := regexp.MatchString(`^--`, tryname)
			if long {
				tryname = tryname[1:]
			}

			if next != "" {
				arguments[next] = "true"
				next = ""
			}
			next = strings.ToUpper(tryname[1:2]) + tryname[2:]
		} else {
			if next != "" {
				arguments[next] = value
				next = ""
			} else {
				arguments["$"+strconv.FormatInt(int64(dolsI), 10)] = value
				dolsI++
			}
		}
	}

	if next != "" {
		arguments[next] = "true"
	}

	arguments["$0"] = command + " " + arguments.Join()

	return command, arguments
}

func getProvider(c string, args Arguments) string {
	gupmConfig := utils.GupmConfig()
	defaultProvider := "gupm"

	if c == "install" {
		defaultProvider = gupmConfig.DefaultProvider
	}

	if defaultProvider == "os" {
		osName := utils.OSNAME()
		if gupmConfig.OsProviders[osName] != "" {
			defaultProvider = gupmConfig.OsProviders[osName]
		} else {
			ui.Error("No provider set for", osName)
			return utils.DIRNAME()
		}
	}

	if utils.FileExists("gupm.json") {
		config, err := utils.ReadGupmJson("gupm.json")
		if err != nil {
			ui.Error(err)
		} else {
			if config.Cli.DefaultProviders[c] != "" {
				defaultProvider = config.Cli.DefaultProviders[c]
			}
		}
	}

	if args["Provider"] != "" {
		ProviderWasForced = true
		return args["Provider"]
	} else if args["P"] != "" {
		ProviderWasForced = true
		return args["P"]
	} else {
		return defaultProvider
	}
}

func ExecCli(c string, args Arguments) (bool, error) {
	var err error
	notFound := "Cannot find commmand"
	shorthands := map[string]string{
		"h":  "help",
		"m":  "make",
		"i":  "install",
		"d":  "delete",
		"p":  "publish",
		"b":  "bootstrap",
		"c":  "cache",
		"s":  "self",
		"t":  "test",
		"pl": "plugin",
	}

	if shorthands[c] != "" {
		c = shorthands[c]
	}

	if provider := getProvider(c, args); provider != "" {
		Provider = provider
	}

	if c == "help" {
		fmt.Println("make / m :", "[--provider=]", "Install projects depdencies based on info in the entry point (depends on provider)")
		fmt.Println("install / i :", "[--provider=]", "Install package")
		fmt.Println("remove / r :", "[--provider=]", "remove package from module config")
		fmt.Println("publish / p :", "[--provider=]", "publish a project based on the model of your specific provider")
		fmt.Println("bootstrap / b :", "[--provider=]", "bootstrap a new project based on the model of your specific provider")
		fmt.Println("test / t :", "[--provider=] Run project's tests in tests folder.")

		fmt.Println("cache / c :", "clear or check the cache with \"cache clear\" or \"cache check\"")
		fmt.Println("self / s :", "self manage gupm. Try g \"self upgrade\" or \"g self uninstall\"")
		fmt.Println("plugin / pl :", "To install a plugin \"g pl install\". Then use \"g pl create\" to create a new one and \"g pl link\" to test your plugin")
	} else if c == "make" {
		BuildGitHooks(".")
		err = InstallProject(".")
	} else if c == "install" {
		err = AddDependency(".", args.AsList())
		if err == nil {
			err = InstallProject(".")
		}
	} else if c == "publish" {
		err = Publish(".", args["$1"])
	} else if c == "delete" {
		err = RemoveDependency(".", args.AsList())
	} else if c == "plugin" {
		if args["$1"] == "create" {
			PluginCreate(".")
		} else if args["$1"] == "link" {
			PluginLink(".")
		} else if args["$1"] == "install" {
			err = PluginInstall(".", args.AsList()[1:])
		} else if args["$1"] == "delete" {
			PluginDelete(".", args.AsList()[1:])
		} else {
			ui.Error(notFound, args["$1"], "\n", "try cache clear or cache check")
		}
	} else if c == "cache" {
		if args["$1"] == "clear" {
			CacheClear()
		} else if args["$1"] == "check" {
			ui.Error("Not implemented yet.")
		} else {
			ui.Error(notFound, args["$1"], "\n", "try cache clear or cache check")
		}
	} else if c == "self" {
		if args["$1"] == "upgrade" {
			SelfUpgrade()
		} else if args["$1"] == "uninstall" {
			SelfUninstall()
		} else {
			ui.Error(notFound, args["$1"])
		}
	} else if c == "bootstrap" {
		err = Bootstrap(".")
	} else if c == "test" {
		RunTest("tests")
	} else if c == "hook" {
		RunHook(".", args["$1"])
	} else {
		return false, nil
	}

	return true, err
}
