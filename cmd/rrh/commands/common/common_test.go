package common

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
)

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
