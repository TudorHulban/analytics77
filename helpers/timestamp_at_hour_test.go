package helpers

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetTodayHourInOffset(t *testing.T) {
	now := time.Date(2026, 7, 3, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		name           string
		hour           int
		utcOffsetHours int
		want           int64
	}{
		{
			name:           "1.UTC midnight hour 0: 2026-07-03 00:00:00 UTC",
			hour:           0,
			utcOffsetHours: 0,
			want:           1783036800,
		},
		{
			name:           "2.UTC hour 12: 2026-07-03 12:00:00 UTC",
			hour:           12,
			utcOffsetHours: 0,
			want:           1783080000,
		},
		{
			name:           "3. UTC+5 hour 10: local 2026-07-03 10:00:00 UTC+5 → UTC 05:00",
			hour:           10,
			utcOffsetHours: 5,
			want:           1783054800,
		},
		{
			name:           "4. UTC-4 hour 8: local 2026-07-03 08:00:00 UTC-4 → UTC 12:00",
			hour:           8,
			utcOffsetHours: -4,
			want:           1783080000,
		},
		{
			name:           "5. UTC+3 hour 0: local 2026-07-03 00:00:00 UTC+3 → UTC 21:00 (previous day)",
			hour:           0,
			utcOffsetHours: 3,
			want:           1783026000,
		},
		{
			name:           "6. UTC-8 hour 23: local 2026-07-03 23:00:00 UTC-8 → UTC 07:00",
			hour:           23,
			utcOffsetHours: -8,
			want:           1783148400,
		},
		{
			name:           "7. Bucharest UTC+3 hour 7: local 2026-07-03 07:00:00 EEST → UTC 04:00",
			hour:           7,
			utcOffsetHours: 3,
			want:           1783051200, // Friday, July 3, 2026 at 4:00:00 AM UTC
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				got := GetHourInOffsetAsUTC(now, tc.hour, tc.utcOffsetHours)

				require.Equal(
					t,
					tc.want,
					got,

					"GetTodayHourInOffset(%d, %d) = %d, want %d difference: %.00f minutes",
					tc.hour,
					tc.utcOffsetHours,
					got,
					tc.want,
					math.Abs(float64(got-tc.want))/60,
				)
			},
		)
	}
}
