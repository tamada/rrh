package common

import (
	"strings"

	"github.com/gookit/color"
)

var labelColorFunc, repoColorFunc, groupColorFunc func(r string) string

var labelColor, repoColor, groupColor string

/*
ColorizedRepositoryID returns the colorized repository id string from configuration.
*/
func ColorizedRepositoryID(repoID string) string {
	return repoColorFunc(repoID)
}

/*
ColorizedGroupName returns the colorized group name string from configuration.
*/
func ColorizedGroupName(groupName string) string {
	return groupColorFunc(groupName)
}

/*
ColorizedLabel returns the colorized label string from configuration.
*/
func ColorizedLabel(label string) string {
	return labelColorFunc(label)
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
	labelColor = ""
	for _, c := range colors {
		parseEach(c)
	}
	updateFuncs()
}

func parseEach(c string) {
	var colors = strings.Split(c, ":")
	switch colors[0] {
	case "repository":
		repoColor = color.ParseCodeFromAttr(strings.Replace(c, "repository:", "", -1))
	case "group":
		groupColor = color.ParseCodeFromAttr(strings.Replace(c, "group:", "", -1))
	case "label":
		labelColor = color.ParseCodeFromAttr(strings.Replace(c, "label:", "", -1))
	}
}

func updateFuncs() {
	repoColorFunc = generateColorFunc(repoColor)
	groupColorFunc = generateColorFunc(groupColor)
	labelColorFunc = generateColorFunc(labelColor)
}

func generateColorFunc(targetColor string) func(s string) string {
	if targetColor != "" {
		var printer = color.NewPrinter(targetColor)
		return func(r string) string {
			return printer.Sprint(r)
		}
	}
	return func(r string) string {
		return r
	}
}

/*
SetColorize sets to enable colorization.
*/
func SetColorize(enable bool) {
	if !enable {
		labelColor = ""
		groupColor = ""
		repoColor = ""
	}
	updateFuncs()
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
