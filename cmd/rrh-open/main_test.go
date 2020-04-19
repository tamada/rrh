package main

import "testing"

func Example_notFound() {
	goMain([]string{"open", "not_exist_repo"})
	// Output:
	// not_exist_repo: repository not found
}

func Example_printHelp() {
	goMain([]string{"open", "-h"})
	// Output:
	// rrh open [OPTIONS] <REPOSITORIES...>
	// OPTIONS
	//     -b, --browser    open the webpage of the specified repository.
	//     -f, --folder     open the folder of the specified repository (Default).
	//     -h, --help       print this message.
	// ARGUMENTS
	//     REPOSITORIES     specifies repository names.
}

func TestConvertGitURL(t *testing.T) {
	testdata := []struct {
		giveString string
		errorFlag  bool
		wontString string
	}{
		{"git@github.com:tamada/rrh.git", false, `https://github.com/tamada/rrh`},
		{"git@github.com:tamada/rrh", false, `https://github.com/tamada/rrh`},
	}
	for _, td := range testdata {
		url, err := convertToRepositoryURL(td.giveString)
		if (err == nil) == td.errorFlag {
			t.Errorf("convertToRepositoryURL(%s) should be %v, but %v", td.giveString, td.errorFlag, err)
		}
		if url != td.wontString {
			t.Errorf("convertToRepositoryURL(%s) wont %s, but got %s", td.giveString, td.wontString, url)
		}
	}
}
