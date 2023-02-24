package main

import (
	"flag"
	"fmt"
	"github.com/go-lean/bevaluate/info"
	"os"
)

func main() {
	target := flag.String("target", "origin/master", "target=origin/production")
	flag.Parse()

	changes, err := info.CollectChanges(*target)
	exitOnError(err, "could not collect changes")

	if len(changes) == 0 {
		return
	}
}

func exitOnError(err error, context string) {
	if err == nil {
		return
	}

	fmt.Printf("%s: %v", context, err)
	os.Exit(1)
}
