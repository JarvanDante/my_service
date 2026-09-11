package logic

import (
	"testing"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

func TestTaggedWorkTable(t *testing.T) {
	cases := []struct {
		ct    int
		table string
		field string
	}{
		{1, "video", "tags"},
		{2, "video", "tags"},
		{3, "video", "tags"},
		{4, "comics", "tags"},
		{5, "photo_album", "tags"},
		{6, "post", "topics"},
		{7, "novel", "tags"},
		{99, "", ""},
	}
	for _, c := range cases {
		table, field := taggedWorkTable(c.ct)
		if table != c.table || field != c.field {
			t.Fatalf("ct=%d got %s.%s", c.ct, table, field)
		}
	}
}

func TestUseCountMapping(t *testing.T) {
	cases := []struct {
		ct    int
		table string
		kind  int
		ok    bool
	}{
		{1, "video", entity.VideoKindVideo, true},
		{2, "video", entity.VideoKindDouyin, true},
		{3, "video", entity.VideoKindCartoon, true},
		{4, "comics", 0, false},
		{5, "video", 0, false},
	}
	for _, c := range cases {
		if supportsUseCount(c.ct) != (c.ct >= 1 && c.ct <= 4) {
			t.Fatalf("ct=%d supports", c.ct)
		}
		if got := useCountTable(c.ct); got != c.table {
			t.Fatalf("ct=%d table %s", c.ct, got)
		}
		kind, ok := useCountKind(c.ct)
		if ok != c.ok || kind != c.kind {
			t.Fatalf("ct=%d kind %d ok=%v", c.ct, kind, ok)
		}
	}
}
