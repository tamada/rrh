package execcmd

import "testing"

func TestValidateArguments(t *testing.T) {
	testdata := []struct {
		execOpts  *execOptions
		args      []string
		wontError bool
	}{
		{&execOptions{}, []string{"ls"}, true},
		{&execOptions{groups: []string{"group1"}}, []string{}, true},
		{&execOptions{repositories: []string{"repo1"}}, []string{"ls"}, false},
	}
	for _, td := range testdata {

		execOpts = td.execOpts
		err := validateArguments(nil, td.args)
		if err == nil && td.wontError || err != nil && !td.wontError {
			t.Errorf("%v wont error %v but got %s", td.args, td.wontError, err.Error())
		}
	}
}
