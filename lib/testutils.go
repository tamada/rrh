package lib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func copyfile(fromfile string) string {
	var content, _ = ioutil.ReadFile(fromfile)
	var file, _ = ioutil.TempFile("../testdata/", "tmp")
	file.Write(content)
	defer file.Close()
	return file.Name()
}

/*
Rollback rollbacks database after executing function f.
*/
func Rollback(dbFile, configFile string, f func(config *Config, db *Database)) string {
	var newDBFile = copyfile(dbFile)
	var newConfigFile = copyfile(configFile)
	defer os.Remove(newConfigFile)
	os.Setenv(RrhConfigPath, newConfigFile)
	os.Setenv(RrhDatabasePath, newDBFile)

	var config = OpenConfig()
	var db, err = Open(config)
	if err != nil {
		fmt.Println(err.Error())
	}

	f(config, db)

	os.Setenv(RrhConfigPath, configFile) // replace the path of config file.
	os.Setenv(RrhDatabasePath, dbFile)

	return newDBFile
}

/*
ReplaceNewline trims spaces and converts the return codes in `originalString` to `replaceTo` string.
*/
func ReplaceNewline(originalString, replaceTo string) string {
	return strings.NewReplacer(
		"\r\n", replaceTo,
		"\r", replaceTo,
		"\n", replaceTo,
	).Replace(strings.TrimSpace(originalString))
}

/*
CaptureStdout is referred from https://qiita.com/kami_zh/items/ff636f15da87dabebe6c.
*/
func CaptureStdout(f func()) string {
	r, w, _ := os.Pipe()
	var stdout = os.Stdout
	os.Stdout = w

	f()

	os.Stdout = stdout
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
