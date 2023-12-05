package util

import "time"

func GetCurrentTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}
