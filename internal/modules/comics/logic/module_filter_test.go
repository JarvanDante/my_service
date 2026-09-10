package logic

import "testing"

func TestParseModuleFilterCompanyJSON(t *testing.T) {
	q := parseModuleFilter(`{"tag_id":"228","order":"rand"}`, nil, nil)
	if len(q.TagIDs) != 1 || q.TagIDs[0] != 228 {
		t.Fatalf("tag %#v", q.TagIDs)
	}
	if q.Order != "rand" {
		t.Fatalf("order %q", q.Order)
	}
	f := q.listFilter(nil, []string{"丝袜"}, 9, false)
	if !f.Shuffle || f.Sort != 2 {
		t.Fatalf("list %#v", f)
	}
}

func TestParseModuleFilterAliases(t *testing.T) {
	q := parseModuleFilter(`{"cat_id":[3,4],"order":"click","is_end":"y","pay_type":"vip","recommend":1,"ids":"10,11"}`, nil, nil)
	if len(q.CatIDs) != 2 || q.CatIDs[0] != 3 || q.CatIDs[1] != 4 {
		t.Fatalf("cat %#v", q.CatIDs)
	}
	if q.Order != "hot" || q.IsEnd != "y" || q.PayType != 1 || !q.Recommend {
		t.Fatalf("q %#v", q)
	}
	if len(q.IDs) != 2 || q.IDs[0] != 10 {
		t.Fatalf("ids %#v", q.IDs)
	}
	f := q.listFilter([]string{"韩漫"}, nil, 6, false)
	if f.Sort != 1 || f.UpdateStatus != 2 || f.PayType != 1 || !f.OnlyRecommend {
		t.Fatalf("list %#v", f)
	}
}

func TestEncodeFilterRoundTrip(t *testing.T) {
	raw := `{"tag_id":"12,13","order":"new","is_end":"n"}`
	got := encodeFilter(parseModuleFilter(raw, nil, nil))
	again := parseModuleFilter(got, nil, nil)
	if len(again.TagIDs) != 2 || again.IsEnd != "n" || again.Order != "new" {
		t.Fatalf("roundtrip %s %#v", got, again)
	}
}

func TestParseCatPosition(t *testing.T) {
	if parseCatPosition("comic_home") != 0 {
		t.Fatal("home")
	}
	if parseCatPosition("cat_29") != 29 {
		t.Fatal("cat")
	}
}
