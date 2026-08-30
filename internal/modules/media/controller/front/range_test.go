package front

import "testing"

func TestParseBytesRange(t *testing.T) {
	cases := []struct {
		h              string
		size           int64
		start, end     int64
		partial        bool
	}{
		{"", 1000, 0, 999, false},
		{"bytes=0-499", 1000, 0, 499, true},
		{"bytes=500-", 1000, 500, 999, true},
		{"bytes=-200", 1000, 800, 999, true},
	}
	for _, c := range cases {
		start, end, partial := parseBytesRange(c.h, c.size)
		if start != c.start || end != c.end || partial != c.partial {
			t.Fatalf("%q: got %d-%d partial=%v want %d-%d partial=%v",
				c.h, start, end, partial, c.start, c.end, c.partial)
		}
	}
}
