package logic

import "testing"

func TestCountTarget(t *testing.T) {
	cases := []struct {
		media, op int
		table     string
		ok        bool
	}{
		{mediaComics, opFav, "comics", true},
		{mediaComics, opLike, "", false},
		{mediaPost, opLike, "post", true},
		{mediaPost, opFav, "", false},
		{mediaNovel, opFav, "novel", true},
		{mediaPhoto, opFav, "photo", true},
		{1, opFav, "", false},
	}
	for _, c := range cases {
		table, ok := countTarget(c.media, c.op)
		if ok != c.ok || table != c.table {
			t.Fatalf("media=%d op=%d got (%s,%v) want (%s,%v)", c.media, c.op, table, ok, c.table, c.ok)
		}
	}
}
