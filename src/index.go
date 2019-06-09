package main

import (
	"os"
	"github.com/spf13/cobra"
	"./utils"
	"./ui"
	"./provider"
	"./jsVm"
	"strings"
	"path/filepath"
	"strconv"
	"regexp"
	"fmt"
	"runtime"
	"time"
)

type json map[string]interface {}

var Provider string

func setProvider(cmd *cobra.Command, args []string) {
	if(Provider == "") {
		Provider = "gupm"
	}
}

var testCmd = &cobra.Command{
	Use:   "test [--provider=]",
	Short: "test a project",
	Long:  `test a project`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunTest("tests", args)
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var tCmd = &cobra.Command{
	Use:   "t [--provider=]",
	Short: "test a new project",
	Long:  `test a new project`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunTest("tests", args)
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap [--provider=]",
	Short: "bootstrap a new project",
	Long:  `bootstrap a new project based on the model of your specific provider`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := Bootstrap(".", args)
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var bCmd = &cobra.Command{
	Use:   "b [--provider=]",
	Short: "bootstrap a new project",
	Long:  `bootstrap a new project based on the model of your specific provider`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := Bootstrap(".", args)
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var publishCmd = &cobra.Command{
	Use:   "publish [--provider=]",
	Short: "publish a project",
	Long:  `publish a project based on the model of your specific provider`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := Publish(".")
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var pCmd = &cobra.Command{
	Use:   "p [--provider=]",
	Short: "publish a project",
	Long:  `publish a project based on the model of your specific provider`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := Publish(".")
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "CLI to help you build plugins",
	Long:  `To install a plugin "g plugin install". Then use "g plugin create" to create a new one and "g plugin link" to test your plugin.`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		if(args[0] == "create") {
			PluginCreate(".")
		} else if (args[0] == "link") {
			PluginLink(".")
		} else if (args[0] == "install") {
			err := PluginInstall(".", args[1:])
			if (err != nil) {
				ui.Error(err.Error())
			}
		} else if (args[0] == "delete") {
			PluginDelete(".", args[1:])
		} else {
			fmt.Println("Unknown command: ", args[0])
		}
	},
}

var plCmd = &cobra.Command{
	Use:   "pl",
	Short: "CLI to help you build plugins",
	Long:  `To install a plugin "g pl install". Then use "g pl create" to create a new one and "g pl link" to test your plugin.`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		if(args[0] == "create") {
			PluginCreate(".")
		} else if (args[0] == "link") {
			PluginLink(".")
		} else if (args[0] == "install") {
			err := PluginInstall(".", args[1:])
			if (err != nil) {
				ui.Error(err.Error())
			}
		} else if (args[0] == "delete") {
			PluginDelete(".", args[1:])
		} else {
			fmt.Println("Unknown command: ", args[0])
		}
	},
}

var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "self manage gupm",
	Long:  `self manage gupm. Try g self upgrade" or "g self uninstall""`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		if(args[0] == "upgrade") {
			SelfUpgrade()
		} else if (args[0] == "uninstall") {
			SelfUninstall()
		} else {
			fmt.Println("Unknown command: ", args[0])
		}
	},
}

var sCmd = &cobra.Command{
	Use:   "s",
	Short: "self manage gupm",
	Long:  `self manage gupm. Try g self upgrade" or "g self uninstall""`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		if(args[0] == "upgrade") {
			SelfUpgrade()
		} else if (args[0] == "uninstall") {
			SelfUninstall()
		} else {
			fmt.Println("Unknown command: ", args[0])
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove package",
	Long:  `remove package from module config`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := RemoveDependency(".", args)
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var rCmd = &cobra.Command{
	Use:   "r",
	Short: "remove package",
	Long:  `remove package from module config`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := RemoveDependency(".", args)
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "cache a new project",
	Long:  `cache a new project based on the model of your specific provider`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		if(args[0] == "clear") {
			CacheClear()
		} else if (args[0] == "clear") {
			ui.Error("Not implemented yet.")
		}
	},
}

var cCmd = &cobra.Command{
	Use:   "c",
	Short: "cache management",
	Long:  `clear or check the cache with "cache clear" or "cache check"`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		if(args[0] == "clear") {
			CacheClear()
		} else if (args[0] == "clear") {
			ui.Error("Not implemented yet.")
		}
	},
}

var installCmd = &cobra.Command{
	Use:   "install [--provider=] package-name",
	Short: "install package",
	Long:  `install package based on info in the entry point (depends on provider)`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := AddDependency(".", args)
		InstallProject(".")
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var iCmd = &cobra.Command{
	Use:   "i [--provider=] package-name",
	Short: "install package",
	Long:  `install package based on info in the entry point (depends on provider)`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := AddDependency(".", args)
		InstallProject(".")
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var makeCmd = &cobra.Command{
	Use:   "make [--provider=]",
	Short: "make package",
	Long:  `make package based on info in the entry point (depends on provider)`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := InstallProject(".")
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}

var mCmd = &cobra.Command{
	Use:   "m [--provider=]",
	Short: "make package",
	Long:  `make package based on info in the entry point (depends on provider)`,
	PreRun: setProvider,
	Run: func(cmd *cobra.Command, args []string) {
		err := InstallProject(".")
		if(err != nil) {
			ui.Error(err.Error())
		} 
	},
}
  
var rootCmd = &cobra.Command{
	Use:   "GuPM",
	Short: "GuPM is the Global Universal Project Manager",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
  
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.Error(err.Error())
	}
}

func ScriptExists(path string) string {
	if (utils.FileExists(path + ".gs")) {
		return path + ".gs"
	} else if(utils.FileExists(path) && !utils.IsDirectory(path)) {
		return path
	} else {
		return ""
	}
}

func executeFile(path string, args []string) {
	i := 1
	next := ""
	input := make(map[string]interface {})
	input["$0"] = strings.Join(args," ")

	if(len(args) > 2) {
		args = args[2:]
	}
	
	for _, value := range args {	
		nameCheck := regexp.MustCompile(`^-(\w+)`)
		tryname := nameCheck.FindString(value)
		if(tryname != "") {
			next = strings.ToUpper(tryname[1:2]) + tryname[2:]
		} else {
			if(next != "") {
				input[next] = value
				next = ""
			} else {
				input["$" + strconv.FormatInt(int64(i), 10)] = value
				i++
			}
		}
	}

	_, err := jsVm.Run(path, input)
	if(err != nil) {
	  ui.Error("File execution failed")
		ui.Error(err.Error())
		os.Exit(1)
	}
}

func binFile(name string, args []string) {
	path := utils.Path("./.bin/"+name)
	realPath, _ := filepath.EvalSymlinks(path)
	utils.ExecCommand(realPath, args)
}

func main() {
	start := time.Now()

	if runtime.GOOS == "darwin" {
		utils.ExecCommand("ulimit", []string{"-n", "2048"})
	}
	
	binFolder := make(map[string]bool)

	if(utils.FileExists(".gupm_rc.gs")) {
		executeFile(".gupm_rc.gs", os.Args)
	}

	if(utils.FileExists(".bin")) {
		for _, file := range utils.ReadDir(".bin") {
			binFolder[file.Name()] = true
		}
	}
	
	packageConfig := new(provider.GupmEntryPoint)
	errConfig := utils.ReadJSON("gupm.json", &packageConfig)
	if(errConfig != nil) {
		fmt.Println("Config file not found.")
	}
	aliases := packageConfig.Cli.Aliases

	rootCmd.AddCommand(makeCmd)
	makeCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(mCmd)
	mCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")

	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(iCmd)
	iCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")

	rootCmd.AddCommand(removeCmd)
	removeCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(rCmd)
	rCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(bCmd)
	bCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	
	rootCmd.AddCommand(publishCmd)
	publishCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(pCmd)
	pCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	
	rootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(tCmd)
	tCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	
	rootCmd.AddCommand(selfCmd)
	rootCmd.AddCommand(sCmd)

	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(cCmd)

	rootCmd.AddCommand(pluginCmd)
	rootCmd.AddCommand(plCmd)

	c := ""
	if(len(os.Args) > 1) {
		c = os.Args[1]
	}

	script := ScriptExists(c)
	if( c == "install" || c == "bootstrap" || c == "make" || c == "update" || c == "cache" ||
	    c == "remove" || c == "self" || c == "plugin" || c == "publish" ||  c == "test" || 
		  c == "i" || c == "b" || c == "m" || c == "u" || c == "c" || c == "r"|| c == "s" || c == "pl" || c == "p" || c == "t") {
			Execute();
			if (script != "") {
				executeFile(script, os.Args)
			}
	} else if (c == "env" || c == "e") {
		toProcess := os.Args[2:]
		re := regexp.MustCompile(`([\w\-\_]+)=([\w\-\_]+)`)
		isEnv := re.FindAllStringSubmatch(toProcess[0], -1)
		for(isEnv != nil) {
			name := isEnv[0][1]
			value := isEnv[0][2]
			os.Setenv(name, value)
			toProcess = toProcess[1:]
			isEnv = re.FindAllStringSubmatch(toProcess[0], -1)
		}
		utils.ExecCommand(toProcess[0], toProcess[1:])
	} else if (aliases[c] != nil) {
		commands := strings.Split(aliases[c].(string), ";")
		for _, command := range commands {
			commandList := strings.Split(command, " ")
			utils.ExecCommand(commandList[0], append(commandList[1:], os.Args[2:]...))
		}
	} else if (binFolder[c] == true) {
		binFile(c, os.Args[2:])	
	} else if (script != "") {
		executeFile(script, os.Args)
	} else if (c == "") {
		fmt.Println("Welcome to GuPM version 1.0.0 \ntry 'g help' for a list of commands. Try 'g filename' to execute a file.")
	} else {
		fmt.Println("Command not found. Try 'g help' or check filename.")
	}

	ui.Stop()
	timeElapsed := fmt.Sprintf("%f", time.Since(start).Seconds())
	fmt.Println(timeElapsed+"s elapsed\n")
}