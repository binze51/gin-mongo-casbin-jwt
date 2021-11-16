package utils

import (
	"fmt"
	"time"
)

func TimeReprOfDuration(duration time.Duration) string {
	return fmt.Sprintf("%.1fms", float32(duration)/float32(time.Millisecond))
}
