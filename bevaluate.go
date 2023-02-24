package main

import (
	"flag"
	"fmt"
	"github.com/go-lean/bevaluate/operations"
	"os"
)

const (
	exitCodeGeneralError = iota + 1
	exitCodeNoCmdSelected
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

	var err error

	cmd := os.Args[1]

	switch cmd {
	case "run":
		err = runCMD.Parse(os.Args[2:])
		exitOnError(err, "could not parse arguments", exitCodeInvalidArgs)

		content := *changes
		if *isFile {
			data, err := os.ReadFile(*changes)
			exitOnError(err, "could not read changes file", exitCodeIOError)

			content = string(data)
		}

		err = operations.EvaluateBuild(root, content)
	default:
		fmt.Printf("unknown cmd: %q\n", cmd)
		os.Exit(1)
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
