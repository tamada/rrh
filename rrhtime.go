package rrh

import (
	"fmt"
	"time"
)

/*
RrhTime represents the time for RRH command, for marshaling a specific format.
*/
type RrhTime struct {
	time time.Time
}

/*
Now returns now time.
*/
func Now() RrhTime {
	return RrhTime{time.Now()}
}

/*
Unix creates and returns the time by specifying the unix time.
*/
func Unix(sec int64, nsec int64) RrhTime {
	return RrhTime{time.Unix(sec, nsec)}
}

func (rt RrhTime) format() string {
	return rt.time.Format("2006-01-02T15:04:05-07:00")
}

/*
UnmarshalJSON is called on unmarshaling JSON.
*/
func (rt *RrhTime) UnmarshalJSON(data []byte) error {
	var t, err = time.Parse("\"2006-01-02T15:04:05-07:00\"", string(data))
	*rt = RrhTime{t}
	return err
}

/*
MarshalJSON is called on marshaling JSON.
*/
func (rt RrhTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, rt.format())), nil
}
