package common

import "strings"

type Resulter interface {
	Err() error
}

type ErrorList []error

func NewErrorList() ErrorList {
	return []error{}
}

func (el ErrorList) Append(err error) ErrorList {
	if el != nil && err != nil {
		el = append(el, err)
	}
	return el
}

func MergeErrors(resulters []Resulter) error {
	errs := ErrorList{}
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

func (el ErrorList) IsErr() bool {
	return el != nil && len(el) > 0
}

func (el ErrorList) IsNil() bool {
	return el == nil || len(el) == 0
}

func (el ErrorList) Error() string {
	messages := []string{}
	for _, err := range el {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "\n")
}
