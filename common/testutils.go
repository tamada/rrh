package common

import (
	"bytes"
	"io"
	"os"
)

/*
Rollback rollback database after executing f.
*/
func Rollback(dbpath string, f func()) {
	os.Setenv(RrhDatabasePath, dbpath)
	var config = OpenConfig()
	var db, _ = Open(config)

	f()

	db.StoreAndClose()
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
