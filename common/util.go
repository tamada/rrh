package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"gopkg.in/src-d/go-git.v4"
)

/*
NormalizePath normalizes given path.

Example:
  given path: /home/username/some/path
  return:     ~/some/path
*/
func NormalizePath(path string) string {
	var home = os.Getenv("HOME")
	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1)
	}
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
FindRemoveUrlFromRepository read remote url of origin from git repository located in given path.
*/
func FindRemoteUrlFromRepository(absPath string) (string, error) {
	var r, err = git.PlainOpen(absPath)
	if err != nil {
		return "", err
	}
	var origin, err2 = r.Remote("origin")
	if err2 != nil {
		return "", err2
	}
	return origin.Config().URLs[0], nil
}
