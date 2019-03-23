package common

import (
	"fmt"
	"time"
)

type RrhTime struct {
	time time.Time
}

func Now() RrhTime {
	return RrhTime{time.Now()}
}

func Unix(sec int64, nsec int64) RrhTime {
	return RrhTime{time.Unix(sec, nsec)}
}

/*
1970-01-01T00:00:00-09:00
*/
func (rt RrhTime) format() string {
	return rt.time.Format("2006-01-02T15:04:05-07:00")
}

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
