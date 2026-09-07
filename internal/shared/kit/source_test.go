package kit

import "testing"

func TestParseSource(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"  DY01  ", "DY01"},
		{"channel://DY01", "DY01"},
		{"channel://dy_01", "dy_01"},
		{"channel://DY01 extra", "DY01"},
		{"share://abc", ""},
		{"agent://abc", ""},
		{"channel://share://x", ""},
		{"official", ""},
		{"404", ""},
		{"invite", ""},
		{"a", ""},
		{"ab", "ab"},
		{"bad!", ""},
		{"中文渠道", ""},
	}
	for _, c := range cases {
		if got := ParseSource(c.in); got != c.want {
			t.Fatalf("ParseSource(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestSourceClipboard(t *testing.T) {
	if got := SourceClipboard("DY01"); got != "channel://DY01" {
		t.Fatalf("got %q", got)
	}
	if got := SourceClipboard("  "); got != "" {
		t.Fatalf("empty should be empty, got %q", got)
	}
}
