package logic

import "testing"

func TestParseChannel(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"  foo  ", "foo"},
		{"channel://promo01", "promo01"},
		{"agent://agent01", "agent01"},
		{"channel://official", ""},
		{"channel://404", ""},
		{"sign=xxx&channel://bar", "bar"},
	}
	for _, c := range cases {
		if got := ParseChannel(c.in); got != c.want {
			t.Fatalf("ParseChannel(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestParseShareParent(t *testing.T) {
	if got := parseShareParent("share://ABC12"); got != "ABC12" {
		t.Fatalf("parseShareParent got %q", got)
	}
}
