package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDayAndHour(t *testing.T) {
	tests := []struct {
		name string

		offsets      TimestampOffsets
		timestampUTC int64

		expectedMonth int8
		expectedDay   int8
		expectedHour  int8
	}{
		{
			name:         "1. Pure UTC - No Offset, Mid-day",
			timestampUTC: 1780932600, // 2026-06-08 15:30:00 UTC
			offsets: TimestampOffsets{
				OffsetUTCHours:     0,
				TimestampDSTSpring: 0,
				TimestampDSTWinter: 0,
			},

			expectedMonth: 6,
			expectedDay:   8,
			expectedHour:  15,
		},
		{
			name:         "2. New York Standard Time (-5h) - No DST active",
			timestampUTC: 1767225600, // 2026-01-01 00:00:00 UTC
			offsets: TimestampOffsets{
				OffsetUTCHours:     -5,         // -5 hours
				TimestampDSTSpring: 1773481200, // March DST (future)
				TimestampDSTWinter: 1761901200, // Nov DST (past)
			},

			expectedMonth: 12,
			expectedDay:   31, // Shipped back to Dec 31, 2025
			expectedHour:  19, // 19:00 PM NY Time
		},
		{
			name:         "3. London Summer Time (+1h DST Active)",
			timestampUTC: 1780932600, // 2026-06-08 15:30:00 UTC
			offsets: TimestampOffsets{
				OffsetUTCHours:     0,          // London standard is 0
				TimestampDSTSpring: 1774755600, // March 29, 2026 01:00:00 UTC
				TimestampDSTWinter: 1792890000, // October 25, 2026 01:00:00 UTC
			},

			expectedMonth: 6,
			expectedDay:   8,
			expectedHour:  16, // 15:30 + 1 hour DST = 16:30
		},
		{
			name:         "4. Exactly on DST Spring Boundary Start",
			timestampUTC: 1774746000, // 2026-03-29 01:00:00 UTC
			offsets: TimestampOffsets{
				OffsetUTCHours:     0,
				TimestampDSTSpring: 1774746000,
				TimestampDSTWinter: 1792890000,
			},

			expectedMonth: 3,
			expectedDay:   29,
			expectedHour:  2,
		},
		{
			name:         "5. Exactly on DST Winter Boundary End (Back to Standard)",
			timestampUTC: 1792890000, // 2026-10-25 01:00:00 UTC
			offsets: TimestampOffsets{
				OffsetUTCHours:     0,
				TimestampDSTSpring: 1774746000,
				TimestampDSTWinter: 1792890000,
			},

			expectedMonth: 10,
			expectedDay:   25,
			expectedHour:  1,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				month, day, hour := ExtractMonthDayHour(
					tc.timestampUTC,
					&tc.offsets,
				)

				assert.EqualValues(t,
					tc.expectedMonth,
					month,

					"E1. computed month different than expected value",
				)

				assert.EqualValues(t,
					tc.expectedDay,
					day,

					"E2. computed day different than expected value",
				)

				assert.EqualValues(t,
					tc.expectedHour,
					hour,

					"E3. computed hour different than expected value",
				)
			},
		)
	}
}
