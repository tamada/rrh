package utils

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
)

func TestJoinArray(t *testing.T) {
	testdata := []struct {
		gives []string
		wont  string
	}{
		{[]string{"1", "2", "3", "4"}, "1, 2, 3 and 4"},
		{[]string{"1", "2", "3"}, "1, 2 and 3"},
		{[]string{"1", "2"}, "1 and 2"},
		{[]string{"1"}, "1"},
		{[]string{}, ""},
	}
	for _, td := range testdata {
		got := JoinArray(td.gives)
		if got != td.wont {
			t.Errorf("JoinArray(%v) did not match, wont %s, got %s", td.gives, td.wont, got)
		}
	}
}

func TestErrorFunc(t *testing.T) {
	err := PerformRrhCommand(nil, []string{}, func(c *cobra.Command, args []string, db *rrh.Database) error {
		return errors.New("some error")
	})
	if err == nil {
		t.Error("some error should be")
	}
}

func TestExamineDB(t *testing.T) {
	PerformRrhCommand(nil, []string{}, func(c *cobra.Command, args []string, db *rrh.Database) error {
		if db == nil {
			t.Error("db was nil")
		}
		if db.Config == nil {
			t.Error("db.Config was nil")
		}
		return nil
	})
}
