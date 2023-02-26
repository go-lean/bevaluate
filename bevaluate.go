package main

import (
	"flag"
	"fmt"
	"github.com/go-lean/bevaluate/config"
	"github.com/go-lean/bevaluate/operations"
	"github.com/go-lean/bevaluate/storage"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	exitCodeGeneralError = iota + 1
	exitCodeNoCmdSelected
	exitCodeConfig
	exitCodeInvalidArgs
	exitCodeIOError
)

func main() {
	runCMD := flag.NewFlagSet("run", flag.ExitOnError)
	changes := runCMD.String("changes", "", `The changes to be processed in --name-status format: either the path to a file or the actual content. e.g. --changes "changes.txt" --file | --changes "$(git diff master --name-status)"`)
	isFile := runCMD.Bool("file", false, `Specifies whether the changes lead to an actual file on disk.`)

	if len(os.Args) < 2 {
		fmt.Println("no cmd selected")
		os.Exit(exitCodeNoCmdSelected)
	}

	root, errWD := os.Getwd()
	exitOnError(errWD, "could not get working directory", exitCodeIOError)

	cfg := config.Default()
	cfgPath := filepath.Join(root, "bevaluate.yaml")

	configFileData, errConfig := os.ReadFile(cfgPath)
	if errConfig == nil {
		errUnmarshal := yaml.Unmarshal(configFileData, &cfg)
		exitOnError(errUnmarshal, "could not unmarshal config file", exitCodeConfig)
	}

	store := storage.Store{}
	var err error
	cmd := os.Args[1]

	switch cmd {
	case "run":
		err = runCMD.Parse(os.Args[2:])
		exitOnError(err, "could not parse arguments", exitCodeInvalidArgs)

		content := *changes
		if *isFile {
			data, errRead := os.ReadFile(*changes)
			exitOnError(errRead, "could not read changes file", exitCodeIOError)

			content = string(data)
		}

		operation := operations.NewEvaluateOperation(store, cfg)
		err = operation.Run(root, content)
	case "init":
		initOperation := operations.NewInitOperation(store)
		err = initOperation.Run(cfgPath)
	default:
		fmt.Printf("unknown cmd: %q\n", cmd)
		os.Exit(exitCodeInvalidArgs)
	}

	exitOnError(err, "could not execute: "+cmd, exitCodeGeneralError)
}

func exitOnError(err error, context string, exitCode int) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %v\n", context, err)
	os.Exit(exitCode)
}
