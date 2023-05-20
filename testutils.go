package rrh

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func copyfile(fromfile string) string {
	var content, err = ioutil.ReadFile(fromfile)
	if err != nil {
		panic(err.Error())
	}
	var file, err2 = ioutil.TempFile(".", "tmp")
	if err2 != nil {
		panic(err.Error())
	}
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
	os.Setenv(ConfigPath, newConfigFile)
	os.Setenv(DatabasePath, newDBFile)

	var config = OpenConfig()
	var db, err = Open(config)
	if err != nil {
		panic(err)
	}

	f(config, db)

	os.Setenv(ConfigPath, configFile) // replace the path of config file.
	os.Setenv(DatabasePath, dbFile)

	return newDBFile
}

func RollbackAlias(dbFile, configFile, aliasFile string, f func(config *Config, db *Database)) string {
	newConfigFile := copyfile(configFile)
	newDBFile := copyfile(dbFile)
	newAliasFile := copyfile(aliasFile)
	defer os.Remove(newConfigFile)
	defer os.Remove(newAliasFile)
	os.Setenv(ConfigPath, newConfigFile)
	os.Setenv(DatabasePath, newDBFile)
	os.Setenv(AliasPath, newAliasFile)

	config := OpenConfig()
	db, err := Open(config)
	if err != nil {
		panic(err)
	}
	f(config, db)
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
GivesStringAsStdin treats given inputFromStdin string as a byte stream from stdin.
*/
func GivesStringAsStdin(inputFromStdin string, f func()) {
	var r, w, _ = os.Pipe()
	var stdin = os.Stdin
	os.Stdin = r
	w.Write([]byte(inputFromStdin))
	w.Close()

	f()

	os.Stdin = stdin
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
