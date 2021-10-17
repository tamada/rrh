package decorator

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/tamada/rrh/common"
)

type colorSettings map[string]string
type colorFuncs map[string](func(r string) string)

func NewColorDecorator(settings string) (*Color, error) {
	c, err := parse(settings)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	c.updateFuncs()
	return c, err
}

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
RepositoryID returns the colorized repository id string from configuration.
*/
func (c *Color) RepositoryID(repoID string) string {
	return c.executeColorFunc("repository", repoID)
}

/*
GroupName returns the colorized group name string from configuration.
*/
func (c *Color) GroupName(groupName string) string {
	return c.executeColorFunc("group", groupName)
}

/*
EnvironmentValue returns the coloried config value from configuration.
*/
func (c *Color) EnvironmentValue(value string) string {
	return c.executeColorFunc("configValue", value)
}

/*
EnvironmentLabel returns the colorized label string from configuration.
*/
func (c *Color) EnvironmentLabel(label string) string {
	return c.executeColorFunc("label", label)
}

func (c *Color) executeColorFunc(label string, value string) string {
	var f, ok = c.funcs[label]
	if ok {
		return f(value)
	}
	return value
}

func parse(colorSettings string) (*Color, error) {
	c := &Color{settings: map[string]string{}, funcs: map[string](func(r string) string){}}
	errs := common.NewErrorList()
	var colors = strings.Split(colorSettings, "+")
	for _, eachColor := range colors {
		err := c.parseEach(eachColor)
		errs = errs.Append(err)
	}
	c.updateFuncs()
	return c, errs.NilOrThis()
}

func (c *Color) parseEach(eachColor string) error {
	var typeAndValue = strings.Split(eachColor, ":")
	if contains(typeAndValue[0], colorLabels) {
		c.settings[typeAndValue[0]] = color.ParseCodeFromAttr(typeAndValue[1])
	} else {
		return fmt.Errorf("%s: unknown color label", typeAndValue[0])
	}
	return nil
}

func contains(value string, values []string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
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
