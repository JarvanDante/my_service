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

func TestMergeNames(t *testing.T) {
	got := MergeNames(NamesCSV("韩漫"), []string{"日漫", "韩漫", " "})
	want := []string{"韩漫", "日漫"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MergeNames = %#v, want %#v", got, want)
	}
}
