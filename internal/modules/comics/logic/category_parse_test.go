package logic

import "testing"

func TestParseCategories(t *testing.T) {
	got := parseCategories("韩漫,日漫，韩漫")
	if len(got) != 2 || got[0] != "韩漫" || got[1] != "日漫" {
		t.Fatalf("got %#v", got)
	}
	if joinCategories([]string{" 日漫", "韩漫", "日漫"}) != "日漫,韩漫" {
		t.Fatalf("join %q", joinCategories([]string{" 日漫", "韩漫", "日漫"}))
	}
}
