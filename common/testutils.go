package common

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
)

var mutex = new(sync.Mutex)

/*
Rollback rollbacks database after executing function f.
*/
func Rollback(dbpath, configPath string, f func()) {
	mutex.Lock()
	os.Setenv(RrhConfigPath, configPath)
	os.Setenv(RrhDatabasePath, dbpath)
	var config = OpenConfig()
	var db, _ = Open(config)
	defer db.StoreAndClose()

	f()
	defer mutex.Unlock()
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
