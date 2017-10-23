package v1

import "time"

func makeNanoTimestamp() int64 {
	return time.Now().UnixNano()
}
