package internal

import (
	"fmt"

	"github.com/tamada/rrh/lib"
)

func printErrors(config *lib.Config, errs []error) int {
	var onError = config.GetValue(lib.RrhOnError)
	if onError != lib.Ignore {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}
	if len(errs) > 0 && (onError == lib.Fail || onError == lib.FailImmediately) {
		return 4
	}
	return 0
}

func isFailImmediately(config *lib.Config) bool {
	return config.GetValue(lib.RrhOnError) == lib.FailImmediately
}
