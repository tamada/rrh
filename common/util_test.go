package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGitRepositoryCheck(t *testing.T) {
	var testcases = []struct {
		path      string
		errorFlag bool
	}{
		{"../testdata/fibonacci", false},
		{"../testdata/database.json", true},
		{"../testdata/other", true},
		{"../not-exist", true},
	}
	for _, testcase := range testcases {
		var absPath, _ = filepath.Abs(testcase.path)
		var err = IsExistAndGitRepository(absPath, testcase.path)
		if (err == nil) == testcase.errorFlag {
			t.Errorf("%s: error wont: %v, got: %v (%v)", testcase.path, testcase.errorFlag, !testcase.errorFlag, err)
		}
	}
}

func TestStrftime(t *testing.T) {
	os.Setenv(RrhTimeFormat, Relative)
	os.Setenv(RrhConfigPath, "../testdata/config.json")

	var now = time.Now()
	var testcases = []struct {
		formatter string
		time      time.Time
		wont      string
	}{
		{Relative, now.Add(time.Minute * -1), "1 minute ago"},
		{Relative, now.Add(time.Hour * -24 * 6), "6 days ago"},
		{Relative, now.Add(time.Hour * -24 * 10), "1 week ago"},
		{Relative, now.Add(time.Hour * -24 * 15), "2 weeks ago"},
		{"2006-01-02 15:04:05", now, now.Format("2006-01-02 15:04:05")},
	}

	var config = OpenConfig()

	for _, test := range testcases {
		os.Setenv(RrhTimeFormat, test.formatter)
		var time = Strftime(test.time, config)
		if time != test.wont {
			t.Errorf("wont: %s, got: %s", test.wont, time)
		}
	}

	os.Unsetenv(RrhTimeFormat)
	os.Unsetenv(RrhConfigPath)
}

func TestRollback(t *testing.T) {
	var file = Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
		var db, _ = Open(OpenConfig())
		db.ForceDeleteGroup("group1")
		db.ForceDeleteGroup("group2")
		db.DeleteRepository("repo1")
		db.DeleteRepository("repo2")
		db.StoreAndClose()
	})
	defer os.Remove(file)

	var db, _ = Open(OpenConfig())
	if !db.HasGroup("group1") || !db.HasGroup("group2") {
		t.Errorf("database did not rollbacked")
	}
	if !db.HasRepository("repo1") || !db.HasRepository("repo2") {
		t.Errorf("database did not rollbacked")
	}
}

func TestReplaceNewline(t *testing.T) {
	var testcases = []struct {
		give      string
		replaceTo string
		wont      string
	}{
		{"a\nb\nc", ",", "a,b,c"},
		{"a\rb\n", ",", "a,b"},
		{"a\nb\rc\r\n", ",", "a,b,c"},
		{"a\nb\rc\r\n", ", ", "a, b, c"},
	}

	for _, tc := range testcases {
		var got = ReplaceNewline(tc.give, tc.replaceTo)
		if got != tc.wont {
			t.Errorf("ReplaceNewLine(%s, %s) wont: %s, got: %s", tc.give, tc.replaceTo, tc.wont, got)
		}
	}
}

func TestCaptureStdout(t *testing.T) {
	var result = CaptureStdout(func() {
		fmt.Println("Hello World")
	})
	result = strings.TrimSpace(result)
	if result != "Hello World" {
		t.Errorf("wont: \"Hello World\", got: %s", result)
	}
}
