package main

import (
	"flag"
	"fmt"
	"github.com/go-lean/bevaluate/info"
	"github.com/go-lean/bevaluate/storage"
	"os"
	"path/filepath"
)

func main() {
	target := flag.String("target", "origin/master", "target=origin/production")
	flag.Parse()

	changes, err := info.CollectChanges(*target)
	exitOnError(err, "could not collect changes")

	if len(changes) == 0 {
		return
	}

	root, err := os.Getwd()
	exitOnError(err, "could not get working directory")

	opener := storage.FileOpener{}
	reader := storage.DirReader{}
	moduleName, err := storage.ReadModuleName(filepath.Join(root, "go.mod"), opener)
	exitOnError(err, "could not read go module name")

	packageReader := info.NewPackageReader(reader, opener)
	packages, err := packageReader.ReadRecursively(root, moduleName)
	exitOnError(err, "could not read packages")

	fmt.Println(packages)
}

func exitOnError(err error, context string) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %v", context, err)
	os.Exit(1)
}
