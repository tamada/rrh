package rrh

import (
	"strings"

	"github.com/gookit/color"
)

type colorSettings map[string]string
type colorFuncs map[string](func(r string) string)

/*
Color struct shows the color settings of RRH.
*/
type Color struct {
	settings colorSettings
	funcs    colorFuncs
}

var colorLabels = []string{
	"repository", "group", "label", "configValue",
}

/*
ColorizedRepositoryID returns the colorized repository id string from configuration.
*/
func (c *Color) ColorizedRepositoryID(repoID string) string {
	return c.executeColorFunc("repository", repoID)
}

/*
ColorizedGroupName returns the colorized group name string from configuration.
*/
func (c *Color) ColorizedGroupName(groupName string) string {
	return c.executeColorFunc("group", groupName)
}

/*
ColorizeConfigValue returns the coloried config value from configuration.
*/
func (c *Color) ColorizeConfigValue(value string) string {
	return c.executeColorFunc("configValue", value)
}

/*
ColorizedLabel returns the colorized label string from configuration.
*/
func (c *Color) ColorizedLabel(label string) string {
	return c.executeColorFunc("label", label)
}

func (c *Color) executeColorFunc(label string, value string) string {
	var f, ok = c.funcs[label]
	if ok {
		return f(value)
	}
	return value
}

/*
ClearColorize clears the color settings.
*/
func (c *Color) ClearColorize() {
	c.parse("")
}

func (c *Color) parse(colorSettings string) {
	for _, label := range colorLabels {
		c.settings[label] = ""
	}
	var colors = strings.Split(colorSettings, "+")
	for _, eachColor := range colors {
		c.parseEach(eachColor)
	}
	c.updateFuncs()
}

func (c *Color) parseEach(eachColor string) {
	var typeAndValue = strings.Split(eachColor, ":")
	if contains(colorLabels, typeAndValue[0]) {
		c.settings[typeAndValue[0]] = color.ParseCodeFromAttr(typeAndValue[1])
	}
}

func (c *Color) updateFuncs() {
	for _, label := range colorLabels {
		var targetColor, ok = c.settings[label]
		if ok {
			c.funcs[label] = generateColorFunc(targetColor)
		} else {
			c.funcs[label] = generateColorFunc("")
		}
	}
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
func (c *Color) SetColorize(enable bool) {
	if !enable {
		for _, label := range colorLabels {
			c.funcs[label] = generateColorFunc("")
		}
	} else {
		c.updateFuncs()
	}
}

/*
InitializeColor is the initialization function of the colorized output.
The function is automatically called on loading the config file.
*/
func InitializeColor(config *Config) *Color {
	var color = Color{colorSettings{}, colorFuncs{}}
	var settingString = config.GetValue(ColorSetting)
	if config.IsSet(EnableColorized) && settingString != "" {
		color.parse(settingString)
	}
	color.updateFuncs()
	return &color
}
