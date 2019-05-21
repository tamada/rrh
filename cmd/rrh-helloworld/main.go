package main

import (
	"fmt"
	"os"

	"github.com/tamada/rrh/common"
)

func main() {
	fmt.Printf("Hello World, %s\n", os.Getenv(common.RrhConfigPath))
}
