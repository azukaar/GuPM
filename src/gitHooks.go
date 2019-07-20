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
		if !utils.FileExists(".git/hooks/pre-push") {
			utils.WriteFile(".git/hooks/pre-push", "g hook prepush")
		}
	}
}

func runhook(hook string) {
	commandList := strings.Split(hook, " ")
	stagedFiles, _ := utils.RunCommand("git", []string{"diff", "--cached", "--name-only"})
	unPushedFiles := ""
	fyg, _ := utils.RunCommand("git", []string{"log", "--branches", "--not", "--remotes", "--name-status", "--oneline"})
	unPushedFilesRaw := strings.Split(fyg, "\n")

	extractFileName := regexp.MustCompile(`^[AM]\b(.*)`)
	for _, v := range unPushedFilesRaw {
		if extracted := extractFileName.FindStringSubmatch(v); len(extracted) > 0 {
			unPushedFiles += strings.Trim(extracted[1], " 	") + "\n"
		}
	}
	unPushedFiles = strings.Trim(unPushedFiles, " ")

	for i, v := range commandList {
		stagedFilesCheck := regexp.MustCompile(`^\$StagedFiles(\(([\w\/\.\*\-\_]+)?\)?)?`)
		tryStagedFiles := stagedFilesCheck.FindStringSubmatch(v)
		unpushedFilesCheck := regexp.MustCompile(`^\$UnpushedFiles(\(([\w\/\.\*\-\_]+)?\)?)?`)
		tryUnpushedFiles := unpushedFilesCheck.FindStringSubmatch(v)

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

		if len(tryUnpushedFiles) == 3 {
			if tryUnpushedFiles[2] == "" {
				commandList[i] = strings.ReplaceAll(unPushedFiles, " ", "\\ ")
				commandList[i] = strings.ReplaceAll(commandList[i], "\n", " ")
			} else {
				commandList[i] = ""
				unpushedFilesList := strings.Split(unPushedFiles, "\n")
				for _, s := range unpushedFilesList {
					isIn, _ := doublestar.Match(tryUnpushedFiles[2], s)
					if isIn {
						commandList[i] += strings.ReplaceAll(s, " ", "\\ ") + " "
					}
				}
				commandList[i] = strings.Trim(commandList[i], " ")
			}
		}
	}

	commandListConsolidated := strings.Join(commandList, " ")
	commandList = strings.Split(commandListConsolidated, " ")

	err := utils.ExecCommand(commandList[0], append(commandList[1:]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runhooklist(hooklist string) {
	hooks := strings.Split(hooklist, ";")
	for _, hook := range hooks {
		runhook(strings.Trim(hook, " "))
	}
}

func runHooks(hooks interface{}) {
	ch := make(chan int)
	listHooks, isArray := hooks.([]interface{})
	if isArray {
		for _, hook := range listHooks {
			go func(hook string) {
				runhooklist(hook)
				ch <- 0
			}(hook.(string))
		}
		for range listHooks {
			<-ch
		}
	} else {
		runhooklist(hooks.(string))
	}
}

func RunHook(path string, hook string) {
	config, _ := utils.ReadGupmJson("gupm.json")

	if hook == "precommit" && config.Git.Hooks.Precommit != nil {
		runHooks(config.Git.Hooks.Precommit)
	}
	if hook == "prepush" && config.Git.Hooks.Prepush != nil {
		runHooks(config.Git.Hooks.Prepush)
	}
}
