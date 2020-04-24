package internal

import (
	"fmt"

	"github.com/tamada/rrh"
)

func printErrors(config *rrh.Config, errs []error) int {
	var onError = config.GetValue(rrh.OnError)
	if onError != rrh.Ignore {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}
	if len(errs) > 0 && (onError == rrh.Fail || onError == rrh.FailImmediately) {
		return 4
	}
	return 0
}

func isFailImmediately(config *rrh.Config) bool {
	return config.GetValue(rrh.OnError) == rrh.FailImmediately
}
