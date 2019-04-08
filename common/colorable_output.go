package common

import (
	"strings"

	"github.com/gookit/color"
)

var repoColorFunc func(r string) string
var groupColorFunc func(r string) string

var repoColor = ""
var groupColor = ""

var supportedForeColor = map[string]color.Color{
	"red":     color.FgRed,
	"black":   color.FgBlack,
	"white":   color.FgWhite,
	"cyan":    color.FgCyan,
	"blue":    color.FgBlue,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"magenta": color.FgMagenta,
}

var supportedBackColor = map[string]color.Color{
	"red":     color.BgRed,
	"black":   color.BgBlack,
	"white":   color.BgWhite,
	"cyan":    color.BgCyan,
	"blue":    color.BgBlue,
	"green":   color.BgGreen,
	"yellow":  color.BgYellow,
	"magenta": color.BgMagenta,
}

/*
ColorrizedRepositoryID returns the colorrized repository id string from configuration.
*/
func ColorrizedRepositoryID(repoID string) string {
	return repoColorFunc(repoID)
}

/*
ColorrizedGroupName returns the colorrized group name string from configuration.
*/
func ColorrizedGroupName(groupName string) string {
	return groupColorFunc(groupName)
}

func parse(colorSettings string) {
	var colors = strings.Split(colorSettings, "+")
	repoColor = ""
	groupColor = ""
	for _, c := range colors {
		if strings.HasPrefix(c, "repository:") {
			repoColor = color.ParseCodeFromAttr(strings.Replace(c, "repository:", "", -1))
		} else if strings.HasPrefix(c, "group:") {
			groupColor = color.ParseCodeFromAttr(strings.Replace(c, "group:", "", -1))
		}
	}
	updateFuncs()
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
	if colorSetting != "" {
		parse(colorSetting)
	}
	updateFuncs()
}
