package rrh

import (
	"bytes"
	"os"
	"testing"
)

func ExampleMessageCenter() {
	var mc = NewMessageCenter()
	mc.PushLog("info level")
	mc.PushVerbose("verbose level")
	mc.Push("warn level", WARN)
	mc.Push("severe level", SEVERE)

	mc.PrintLog(os.Stdout)
	// Output:
	// info level
	// warn level
	// severe level
}

func TestMessageCenter(t *testing.T) {
	var mc = NewMessageCenter()
	mc.PushLog("info level")
	mc.PushVerbose("verbose level")
	mc.Push("warn level", WARN)
	mc.Push("severe level", SEVERE)
	var buffer = new(bytes.Buffer)
	var wontString = `info level
verbose level
warn level
severe level
`

	mc.PrintVerbose(buffer)
	if string(buffer.Bytes()) != wontString {
		t.Errorf("output string did not match, wont %s, got %s", wontString, string(buffer.Bytes()))
	}
}
