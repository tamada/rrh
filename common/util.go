package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

/*
NormalizePath normalizes given path.

Example:
  given path: /home/username/some/path
  return:     ~/some/path
*/
func NormalizePath(path string) string {
	// var home = os.Getenv("HOME")
	// if strings.HasPrefix(path, home) {
	// 	return strings.Replace(path, home, "~", 1)
	// }
	return path
}

func IsInputYes(prompt string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var text = strings.ToLower(scanner.Text())
	return text == "yes" || text == "y"
}

func ToAbsolutePath(path string, config *Config) string {
	var home = os.Getenv("HOME")
	return strings.Replace(path, "~", home, 1)
}

func CreateParentDir(path string) error {
	var parent = filepath.Dir(path)
	return os.MkdirAll(parent, 0755)
}

func Strftime(before time.Time, config *Config) string {
	var format = config.GetValue(RrhTimeFormat)
	if format != Relative {
		return before.Format(format)
	}
	return humanize.Time(before)
}

/*
CaptureStdout is referred from https://qiita.com/kami_zh/items/ff636f15da87dabebe6c.
*/
func CaptureStdout(f func()) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	var stdout = os.Stdout
	os.Stdout = w

	f()

	os.Stdout = stdout
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), nil
}
