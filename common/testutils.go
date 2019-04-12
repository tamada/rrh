package common

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func copyfile(fromfile string) string {
	var content, _ = ioutil.ReadFile(fromfile)
	var file, _ = ioutil.TempFile("../testdata/", "tmp")
	file.Write(content)
	return file.Name()
}

/*
WithDatabase introduce mutex for using database for only one routine at once.
*/
func WithDatabase(dbFile, configFile string, f func()) string {
	var newDBFile = copyfile(dbFile)
	var newConfigFile = copyfile(configFile)
	os.Setenv(RrhConfigPath, newConfigFile)
	os.Setenv(RrhDatabasePath, newDBFile)

	f()

	defer os.Setenv(RrhConfigPath, configFile) // replace the path of config file.
	defer os.Remove(newConfigFile)
	return newDBFile
}

/*
Rollback rollbacks database after executing function f.
*/
func Rollback(dbpath, configPath string, f func()) string {
	return WithDatabase(dbpath, configPath, func() {
		var config = OpenConfig()
		var db, _ = Open(config)
		defer db.StoreAndClose()

		f()
	})
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
