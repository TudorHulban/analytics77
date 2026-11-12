package fixtures

import (
	"time"

	"github.com/tudorhulban/analytics77/domain/offsets"
)

func DefaultOffset(customTime ...time.Time) offsets.TimestampOffsets {
	var now time.Time

	if len(customTime) == 1 {
		now = customTime[0]
	} else {
		now = time.Now()
	}

	_, offsetSeconds := now.Zone()

	return offsets.TimestampOffsets{
		OffsetUTCHours: int64(offsetSeconds / 3600),
	}
}
