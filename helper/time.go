package helper

import (
	"fmt"
	"time"
)

func FormatDurationToMs(d time.Duration) string {
	return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/float64(time.Millisecond))
}
