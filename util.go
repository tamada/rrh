package rrh

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

/*
IsInputYes print the given prompt and returns TRUE if the user inputs "yes".
*/
func IsInputYes(prompt string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var text = strings.ToLower(scanner.Text())
	return text == "yes" || text == "y"
}

/*
CreateParentDir create the parent directories of the given path.
*/
func CreateParentDir(path string) error {
	var parent = filepath.Dir(path)
	return os.MkdirAll(parent, 0755)
}

/*
Strftime returns the string of the given time.
*/
func Strftime(time time.Time, config *Config) string {
	var format = config.GetValue(TimeFormat)
	if format != Relative {
		return time.Format(format)
	}
	return HumanizeTime(time)
}

/*
HumanizeTime convert the given time to human friendly formatted string.
*/
func HumanizeTime(time time.Time) string {
	return humanize.Time(time)
}

/*
IsExistAndGitRepository checks the given absPath is exist and shows the git repository.
*/
func IsExistAndGitRepository(absPath string, path string) error {
	var fmode, err = os.Stat(absPath)
	if err != nil {
		return err
	}
	if !fmode.IsDir() {
		return fmt.Errorf("%s: not directory", path)
	}
	fmode, err = os.Stat(filepath.Join(absPath, ".git"))
	// If the repository of path is submodule, `.git` will be a file to indicate the `.git` directory.
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: not git repository", path)
	}
	return nil
}

func LoadJson(filePath string, v interface{}) error {
	reader, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func StoreJson(filePath string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	writer, err := os.OpenFile(filePath, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	n, err := writer.Write(bytes)
	if n != len(bytes) {
		return fmt.Errorf("cannot write fully %d/%d", n, len(bytes))
	}
	return err
}

func FindIn(target string, list []string) bool {
	for _, item := range list {
		if target == item {
			return true
		}
	}
	return false
}
