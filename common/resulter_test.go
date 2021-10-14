package common

import (
	"errors"
	"testing"
)

func addNil(el ErrorList) ErrorList {
	return el.Append(nil)
}

func addSome(el ErrorList) ErrorList {
	return el.Append(errors.New("some"))
}

func addEmptyList(el ErrorList) ErrorList {
	return el.Append(NewErrorList())
}

func addSomeList(el ErrorList) ErrorList {
	el2 := addSome(NewErrorList())
	return el.Append(el2)
}

func TestErrorList(t *testing.T) {
	testdata := []struct {
		before        func(el ErrorList) ErrorList
		nilFlag       bool
		nilOrThisFlag bool
	}{
		{func(el ErrorList) ErrorList { return el }, true, true},
		{addNil, true, true},
		{addSome, false, false},
		{addEmptyList, true, true},
		{addSomeList, false, false},
	}
	for _, td := range testdata {
		el := NewErrorList()
		el = td.before(el)
		if el.IsNil() != td.nilFlag {
			t.Errorf("IsNil wont %v, but got %v", !td.nilFlag, el.IsNil())
		}
		if el.IsErr() == td.nilFlag {
			t.Errorf("IsErr wont %v, but got %v", td.nilFlag, el.IsErr())
		}
		err := el.NilOrThis()
		if (err != nil) == td.nilOrThisFlag {
			t.Errorf("NilOrThis wont %v, but got %v", td.nilOrThisFlag, err)
		}
	}
}
