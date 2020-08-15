package trafikverket

import (
	"testing"
	"time"
)

func Test_duration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		dur      time.Duration
		expected string
	}{
		{"now", 0, "$now"},
		{"one day / 24 hours", 24 * time.Hour, "$dateadd(1.0:00:00)"},
		{"two hours", 2 * time.Hour, "$dateadd(0.2:00:00)"},
		{"one minute", time.Minute, "$dateadd(0.0:01:00)"},
		{"one second", time.Second, "$dateadd(0.0:00:01)"},
		{"two days ten hours fifteen minutes", 58*time.Hour + 15*time.Minute, "$dateadd(2.10:15:00)"},
		{"59 minutes", 59 * time.Minute, "$dateadd(0.0:59:00)"},
		{"60 minutes", 60 * time.Minute, "$dateadd(0.1:00:00)"},
		{"23 hours", 23 * time.Hour, "$dateadd(0.23:00:00)"},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := duration(tc.dur); got != tc.expected {
				t.Errorf("duration(%s) = %v, want %v", tc.dur.String(), got, tc.expected)
			}
		})
	}
}
