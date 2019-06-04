package main

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/mitchellh/go-homedir"
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

var cacheCmd = &cobra.Command{
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

var cCmd = &cobra.Command{
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

var bootstrapCmd = &cobra.Command{
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

var bCmd = &cobra.Command{
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
	} else if(utils.FileExists(path)) {
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
	path := "./.bin/"+name
	realPath, _ := filepath.EvalSymlinks(path)
	utils.ExecCommand(realPath, args)
}

func main() {
	start := time.Now()

	hdir, errH := homedir.Dir()
	if(errH != nil) {
		fmt.Println(errH)
		hdir = "."
	}
	flog, _ := os.OpenFile(hdir + "/.gupm/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	os.Stderr = flog

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
	
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(bCmd)
	bCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")

	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(cCmd)

	c := ""
	if(len(os.Args) > 1) {
		c = os.Args[1]
	}

	script := ScriptExists(c)
	if( c == "install" || c == "bootstrap" || c == "make" || c == "uninstall" || c == "cache" ||
		c == "i" || c == "b" || c == "m" || c == "u" || c == "c") {
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
		utils.ExecCommand(aliases[c].(string), os.Args[2:])	
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