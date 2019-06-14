package internal

import (
	"testing"
	"time"

	"github.com/tamada/rrh/lib"
)

func TestStrftime(t *testing.T) {
	var statusOptions = &statusOptions{}
	var targetTime = time.Now().Add(time.Hour * 24 * -7)
	var testcases = []struct {
		timeFormat  string
		time        *time.Time
		wontMessage string
	}{
		{lib.Relative, &targetTime, "1 week ago"},
		{lib.Relative, nil, ""},
		{notSpecified, &targetTime, "1 week ago"},
		{absolute, &targetTime, targetTime.Format("2006-01-02 03:04:05-07")},
	}
	var config = lib.OpenConfig()
	for _, tc := range testcases {
		statusOptions.format = tc.timeFormat

		var string = statusOptions.strftime(tc.time, config)
		if string != tc.wontMessage {
			t.Errorf("%s did not match, wont %s, got %s", tc.timeFormat, tc.wontMessage, string)
		}
	}
}

func TestHelpOfStatusCommand(t *testing.T) {
	var helpMessage = `rrh status [OPTIONS] [REPOSITORIES|GROUPS...]
OPTIONS
    -b, --branches               show the status of the local branches.
    -r, --remote                 show the status of the remote branches.
    -c, --csv                    print result in csv format.
    -f, --time-format <FORMAT>   specifies time format. Available value is
                                 'relative' ad 'absolute'
ARGUMENTS
    REPOSITORIES                 target repositories.  If no repository was specified
                                 the command shows the result of the default group.
    GROUPS                       target groups.  If no group was specified,
                                 the command shows the result of the default group.`
	var command, _ = StatusCommandFactory()
	if command.Help() != helpMessage {
		t.Errorf("help message did not match")
	}
}
