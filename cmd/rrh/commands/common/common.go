package common

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"gopkg.in/go-playground/validator.v9"
)

func PerformRrhCommand(c *cobra.Command, args []string, f func(c *cobra.Command, args []string, db *rrh.Database) error) error {
	var config = rrh.OpenConfig()
	var db, err = rrh.Open(config)
	if err != nil {
		return err
	}
	return f(c, args, db)
}

var structValidator = validator.New()

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
