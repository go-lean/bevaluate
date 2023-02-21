package main

import (
	"flag"
	"fmt"
)

func main() {
	target := flag.String("target", "origin/master", "target=origin/production")
	flag.Parse()

	fmt.Println("hello from bevaluate!")
	fmt.Println("target is " + *target)
}
