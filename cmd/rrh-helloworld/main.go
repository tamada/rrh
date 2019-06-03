package main

import (
	"fmt"
	"os"

	"github.com/tamada/rrh/lib"
)

func main() {
	fmt.Printf("Hello World, %s\n", os.Getenv(lib.RrhConfigPath))
}
