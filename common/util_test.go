package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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

func TestCaptureStdout(t *testing.T) {
	var result, _ = CaptureStdout(func() {
		fmt.Println("Hello World")
	})
	result = strings.TrimSpace(result)
	if result != "Hello World" {
		t.Errorf("wont: \"Hello World\", got: %s", result)
	}
}

func TestNormalizePath(t *testing.T) {
	var testcases = []struct {
		input string
		wont  string
	}{
		{"path1", "path1"},
		{"path2", "path2"},
		{"path3", "path3"},
	}

	for _, test := range testcases {
		if val := NormalizePath(test.input); val != test.wont {
			t.Errorf("wont: %s, got: %s", test.wont, val)
		}
	}
}

func TestToAbsolutePath(t *testing.T) {
	var home = os.Getenv("HOME")
	var testcases = []struct {
		input string
		wont  string
	}{
		{"~/go/src/github.com/tamada/rrh", filepath.Join(home, "go/src/github.com/tamada/rrh")},
		{"go/src/github.com/tamada/rrh", "go/src/github.com/tamada/rrh"},
	}

	os.Setenv(RrhConfigPath, "../testdata/config.json")
	var config = OpenConfig()
	for _, test := range testcases {
		if val := ToAbsolutePath(test.input, config); val != test.wont {
			t.Errorf("wont: %s, got: %s", test.wont, val)
		}
	}
}
