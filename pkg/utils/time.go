package xutils

import "time"

func GetTimeNow() time.Time {
	return time.Now().Local()
}
