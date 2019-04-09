package common

import (
	"strings"

	"github.com/gookit/color"
)

var repoColorFunc func(r string) string
var groupColorFunc func(r string) string

var repoColor = ""
var groupColor = ""

/*
ColorizedRepositoryID returns the colorrized repository id string from configuration.
*/
func ColorizedRepositoryID(repoID string) string {
	return repoColorFunc(repoID)
}

/*
ColorizedGroupName returns the colorrized group name string from configuration.
*/
func ColorizedGroupName(groupName string) string {
	return groupColorFunc(groupName)
}

/*
ClearColorize clears the color settings.
*/
func ClearColorize() {
	parse("")
}

func parse(colorSettings string) {
	var colors = strings.Split(colorSettings, "+")
	repoColor = ""
	groupColor = ""
	for _, c := range colors {
		parseEach(c)
	}
	updateFuncs()
}

func parseEach(c string) {
	if strings.HasPrefix(c, "repository:") {
		repoColor = color.ParseCodeFromAttr(strings.Replace(c, "repository:", "", -1))
	} else if strings.HasPrefix(c, "group:") {
		groupColor = color.ParseCodeFromAttr(strings.Replace(c, "group:", "", -1))
	}
}

func updateFuncs() {
	updateRepoFunc(repoColor)
	updateGroupFunc(groupColor)
}

func updateRepoFunc(repoColor string) {
	if repoColor != "" {
		var printer = color.NewPrinter(repoColor)
		repoColorFunc = func(r string) string {
			return printer.Sprint(r)
		}
	} else {
		repoColorFunc = func(r string) string {
			return r
		}
	}
}

func updateGroupFunc(groupColor string) {
	if groupColor != "" {
		var printer = color.NewPrinter(groupColor)
		groupColorFunc = func(r string) string {
			return printer.Sprint(r)
		}
	} else {
		groupColorFunc = func(r string) string {
			return r
		}
	}
}

/*
InitializeColor is the initialization function of the colorized output.
The function is automatically called on loading the config file.
*/
func InitializeColor(config *Config) {
	var colorSetting = config.GetValue(RrhColor)
	if config.IsSet(RrhEnableColorized) && colorSetting != "" {
		parse(colorSetting)
	}
	updateFuncs()
}
