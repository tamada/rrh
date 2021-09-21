package common

import "strings"

type Resulter interface {
	Err() error
}

type errorList []error

func MergeErrors(resulters []Resulter) error {
	errs := errorList{}
	for _, resulter := range resulters {
		err := resulter.Err()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (el errorList) Error() string {
	messages := []string{}
	for _, err := range el {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, ",")
}
