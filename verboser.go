package rrh

import "github.com/spf13/viper"

type Verboser interface {
	Print(i ...interface{})
	PrintErr(i ...interface{})
	PrintErrf(format string, i ...interface{})
	PrintErrln(i ...interface{})
	Printf(format string, i ...interface{})
	Println(i ...interface{})
}

func PrintIfVerbose(verboser Verboser, message string) {
	if viper.GetBool("verbose") {
		verboser.Println(message)
	}
}
