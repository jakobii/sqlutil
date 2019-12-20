package sqlutil

import "time"

func Date (t time.Time) string {
	return t.Format("2006-01-02")
}
