package main

import (
	"flag"
	"fmt"
	"github.com/go-lean/bevaluate/operations"
	"os"
)

func main() {
	runCMD := flag.NewFlagSet("run", flag.ExitOnError)
	target := runCMD.String("target", "origin/master", "target=origin/production")

	initCMD := flag.NewFlagSet("init", flag.ExitOnError)
	defaultSettings := initCMD.Bool("default", false, "init -default")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("no cmd selected")
		os.Exit(1)
	}

	root, errWD := os.Getwd()
	exitOnError(errWD, "could not get working directory")

	var err error

	cmd := os.Args[1]

	switch cmd {
	case "run":
		err = operations.EvaluateBuild(root, *target)
	case "init":
		if *defaultSettings {
			err = operations.InitDefault()
			break
		}
		err = operations.Init()
	default:
		fmt.Printf("unknown cmd: %q\n", cmd)
		os.Exit(1)
	}

	exitOnError(err, "could not execute: "+cmd)
}

func exitOnError(err error, context string) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %v", context, err)
	os.Exit(1)
}
