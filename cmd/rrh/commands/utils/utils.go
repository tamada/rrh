package utils

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"gopkg.in/go-playground/validator.v9"
)

func PerformRrhCommand(c *cobra.Command, args []string, f func(c *cobra.Command, args []string, db *rrh.Database) error) error {
	if c != nil {
		c.SilenceUsage = true
	}
	config := rrh.OpenConfig()
	db, err := rrh.Open(config)
	if err != nil {
		return err
	}
	return f(c, args, db)
}

func IsVerbose(c *cobra.Command) bool {
	flag, err := c.PersistentFlags().GetBool("verbose")
	return err != nil && flag
}

var structValidator = validator.New()

func ValidateValue(value string, availables []string) error {
	return ValidateValues([]string{value}, availables)
}

func ValidateValues(values []string, availables []string) error {
	no := []string{}
	for _, value := range values {
		lowerValue := strings.ToLower(value)
		if !rrh.FindIn(lowerValue, availables) {
			no = append(no, value)
		}
	}
	if len(no) == 0 {
		return nil
	} else if len(no) == 1 {
		return fmt.Errorf("%v: not available entry. availables: %v", JoinArray(no), JoinArray(availables))
	}
	return fmt.Errorf("%v: not available entries. availables: %v", JoinArray(no), JoinArray(availables))
}

func JoinArray(array []string) string {
	switch len(array) {
	case 0:
		return ""
	case 1:
		return array[0]
	case 2:
		return array[0] + " and " + array[1]
	default:
		newArray := []string{array[0] + ", " + array[1]}
		newArray = append(newArray, array[2:]...)
		return JoinArray(newArray)
	}
}

func ValidateOptions(s interface{}) error {
	errs := structValidator.Struct(s)
	return extractValidationErros(errs)
}

func extractValidationErros(err error) error {
	if err != nil {
		errorText := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errorText = append(errorText, validationErrorToText(err))
		}
		return fmt.Errorf("parameter errors: %s", strings.Join(errorText, "\n\t"))
	}
	return nil
}

func validationErrorToText(e validator.FieldError) string {
	f := e.Field()
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", f)
	case "max":
		return fmt.Sprintf("%s cannot be greater than %s", f, e.Param())
	case "min":
		return fmt.Sprintf("%s must be greater than %s", f, e.Param())
	}
	return fmt.Sprintf("%s is invalid %s", e.Field(), e.Value())
}
