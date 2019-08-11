package server

import (
	"time"
)

func NowUnix() int64 {
	return time.Now().UnixNano() / 1e6
}
