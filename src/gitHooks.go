package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"./utils"
	"github.com/bmatcuk/doublestar"
)

func BuildGitHooks(path string) {
	if utils.FileExists(".git") {
		if !utils.FileExists(".git/hooks/pre-commit") {
			utils.WriteFile(".git/hooks/pre-commit", "g hook precommit")
		}
		if !utils.FileExists(".git/hooks/pre-publish") {
			utils.WriteFile(".git/hooks/pre-publish", "g hook prepublish")
		}
	}
}

func runhook(hook string) {
	commandList := strings.Split(hook, " ")
	stagedFiles, _ := utils.RunCommand("git", []string{"diff", "--cached", "--name-only"})

	for i, v := range commandList {
		stagedFilesCheck := regexp.MustCompile(`^\$StagedFiles(\(([\w\/\.\*\-\_]+)?\)?)?`)
		tryStagedFiles := stagedFilesCheck.FindStringSubmatch(v)

		if len(tryStagedFiles) == 3 {
			if tryStagedFiles[2] == "" {
				commandList[i] = strings.ReplaceAll(stagedFiles, " ", "\\ ")
				commandList[i] = strings.ReplaceAll(commandList[i], "\n", " ")
			} else {
				commandList[i] = ""
				stagedFilesList := strings.Split(stagedFiles, "\n")
				for _, s := range stagedFilesList {
					isIn, _ := doublestar.Match(tryStagedFiles[2], s)
					if isIn {
						commandList[i] += strings.ReplaceAll(s, " ", "\\ ") + " "
					}
				}
				commandList[i] = strings.Trim(commandList[i], " ")
			}
		}
	}
	err := utils.ExecCommand(commandList[0], append(commandList[1:]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RunHook(path string, hook string) {
	config, _ := utils.ReadGupmJson("gupm.json")

	if hook == "precommit" && config.Git.Hooks.Precommit != "" {
		runhook(config.Git.Hooks.Precommit)
	}
	if hook == "prepush" && config.Git.Hooks.Prepush != "" {
		runhook(config.Git.Hooks.Prepush)
	}
}
