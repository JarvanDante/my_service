package kit

import (
	"reflect"
	"testing"
)

func TestNamesCSV(t *testing.T) {
	got := NamesCSV(" 韩漫, 日漫，韩漫 ,")
	want := []string{"韩漫", "日漫"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NamesCSV = %#v, want %#v", got, want)
	}
}

func TestReplaceNameCSV(t *testing.T) {
	if got := ReplaceNameCSV("标签一,标签二", "标签一", "标签甲"); got != "标签甲,标签二" {
		t.Fatalf("got %q", got)
	}
	if got := ReplaceNameCSV("标签一， 标签一", "标签一", "新名"); got != "新名" {
		t.Fatalf("dedupe got %q", got)
	}
	raw := "标签一,标签二"
	if got := ReplaceNameCSV(raw, "没有", "新"); got != raw {
		t.Fatalf("miss got %q", got)
	}
}

func TestMergeNames(t *testing.T) {
	got := MergeNames(NamesCSV("韩漫"), []string{"日漫", "韩漫", " "})
	want := []string{"韩漫", "日漫"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MergeNames = %#v, want %#v", got, want)
	}
}
