package main

import (
	"fmt"
	"os"
)

func goMain(args []string) int {
	var world = "World"
	if len(args) >= 2 {
		world = args[1]
	}
	fmt.Printf("Hello, %s!\n", world)
	return 0
}

func main() {
	var status = goMain(os.Args)
	os.Exit(status)
}
