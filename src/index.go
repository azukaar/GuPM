package main

import (
	"os"
	"github.com/spf13/cobra"
	"./utils"
	"./jsVm"
	"strings"
	"path/filepath"
	"strconv"
	"os/exec"
	"regexp"
	"fmt"
	"time"
)

type json map[string]interface {}

var Provider string

var installCmd = &cobra.Command{
	Use:   "install [--provider=] package-name",
	Short: "install package",
	Long:  `install package based on info in the entry point (depends on provider)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		err := AddDependency(".", args)
		if(err != nil) {
			fmt.Println(err)
		} 
	},
}

var iCmd = &cobra.Command{
	Use:   "i [--provider=] package-name",
	Short: "install package",
	Long:  `install package based on info in the entry point (depends on provider)`,
	Run: func(cmd *cobra.Command, args []string) {
		err := AddDependency(".", args)
		if(err != nil) {
			fmt.Println(err)
		} 
	},
}

var makeCmd = &cobra.Command{
	Use:   "make [--provider=]",
	Short: "make package",
	Long:  `make package based on info in the entry point (depends on provider)`,
	Run: func(cmd *cobra.Command, args []string) {
		err := InstallProject(".")
		if(err != nil) {
			fmt.Println(err)
		} 
	},
}

var mCmd = &cobra.Command{
	Use:   "m [--provider=]",
	Short: "make package",
	Long:  `make package based on info in the entry point (depends on provider)`,
	Run: func(cmd *cobra.Command, args []string) {
		err := InstallProject(".")
		if(err != nil) {
			fmt.Println(err)
		} 
	},
}
  
var rootCmd = &cobra.Command{
	Use:   "GuPM",
	Short: "GuPM is the Global Universal Project Manager",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test")
	},
}
  
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ScriptExists(path string) string {
	if(utils.FileExists(path)) {
		return path
	} else if (utils.FileExists(path + ".gs")) {
		return path + ".gs"
	} else {
		return ""
	}
}

func executeFile(path string, args []string) {
	i := 1
	next := ""
	input := make(map[string]interface {})
	input["$0"] = strings.Join(args," ")
	
	for _, value := range args[2:] {	
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
		fmt.Println("File execution failed")
		fmt.Println(err)
		os.Exit(1)
	}
}

func binFile(name string, args []string) {
	path := "./.bin/"+name
	realPath, _ := filepath.EvalSymlinks(path)
	bashargs := []string{"-c"}
	bashargs = append(bashargs, realPath)
	bashargs = append(bashargs, args...)
	
	cmd := exec.Command("/bin/bash", bashargs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}

func main() {
	start := time.Now()

	binFolder := make(map[string]bool)
	for _, file := range utils.ReadDir("./.bin") {
		binFolder[file.Name()] = true
	}

	rootCmd.AddCommand(makeCmd)
	makeCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(mCmd)
	mCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")

	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")
	rootCmd.AddCommand(iCmd)
	iCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")

	c := os.Args[1]
	script := ScriptExists(c)
	if( c == "install" || c == "make" || c == "uninstall" ||
		c == "i" || c == "m" || c == "u") {
			Execute();
			if (script != "") {
				executeFile(script, os.Args)
			}
	} else if (binFolder[c] == true) {
		binFile(c, os.Args[2:])	
	} else if (script != "") {
		executeFile(script, os.Args)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}