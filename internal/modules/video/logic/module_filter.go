package logic

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/JarvanDante/my_service/internal/modules/video/service"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

type moduleQuery struct {
	TagIDs   []int64
	CatIDs   []int64
	Order    string
	IDs      []int64
	Keywords string
}

func parseCatPosition(pos string) int64 {
	pos = strings.TrimSpace(pos)
	if !strings.HasPrefix(pos, "cat_") {
		return 0
	}
	id, _ := strconv.ParseInt(strings.TrimPrefix(pos, "cat_"), 10, 64)
	if id < 0 {
		return 0
	}
	return id
}

func parseIDList(raw json.RawMessage) []int64 {
	if len(raw) == 0 {
		return nil
	}
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return nil
	}
	if strings.HasPrefix(s, "[") {
		var arr []int64
		if json.Unmarshal(raw, &arr) == nil {
			return arr
		}
		var strs []string
		if json.Unmarshal(raw, &strs) == nil {
			return kit.ParseI64CSV(strings.Join(strs, ","))
		}
		return nil
	}
	if strings.HasPrefix(s, `"`) {
		var str string
		if json.Unmarshal(raw, &str) != nil {
			return nil
		}
		return kit.ParseI64CSV(str)
	}
	var n int64
	if json.Unmarshal(raw, &n) == nil && n > 0 {
		return []int64{n}
	}
	return kit.ParseI64CSV(strings.Trim(s, `"'`))
}

func normalizeOrder(order string) string {
	switch strings.ToLower(strings.TrimSpace(order)) {
	case "rand", "random":
		return "rand"
	case "hot", "click", "ranking", "like":
		return "hot"
	case "new", "update_date", "":
		return "new"
	default:
		return "new"
	}
}

func encodeIDCSV(ids []int64) string {
	if len(ids) == 0 {
		return ""
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			parts = append(parts, strconv.FormatInt(id, 10))
		}
	}
	return strings.Join(parts, ",")
}

func encodeFilter(q moduleQuery) string {
	m := map[string]any{"order": normalizeOrder(q.Order)}
	if s := encodeIDCSV(q.TagIDs); s != "" {
		m["tag_id"] = s
	}
	if s := encodeIDCSV(q.CatIDs); s != "" {
		m["cat_id"] = s
	}
	if s := encodeIDCSV(q.IDs); s != "" {
		m["ids"] = s
	}
	if kw := strings.TrimSpace(q.Keywords); kw != "" {
		m["keywords"] = kw
	}
	b, err := json.Marshal(m)
	if err != nil {
		return `{"order":"new"}`
	}
	return string(b)
}

func parseModuleFilter(raw string, fallbackCats, fallbackTags []int64) moduleQuery {
	q := moduleQuery{Order: "new", CatIDs: append([]int64{}, fallbackCats...), TagIDs: append([]int64{}, fallbackTags...)}
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "{}" || raw == "null" {
		return q
	}
	var in struct {
		TagID    json.RawMessage `json:"tag_id"`
		CatID    json.RawMessage `json:"cat_id"`
		Order    string          `json:"order"`
		IDs      json.RawMessage `json:"ids"`
		Keywords string          `json:"keywords"`
	}
	if json.Unmarshal([]byte(raw), &in) != nil {
		return q
	}
	if ids := parseIDList(in.TagID); len(ids) > 0 {
		q.TagIDs = ids
	}
	if ids := parseIDList(in.CatID); len(ids) > 0 {
		q.CatIDs = ids
	}
	q.Order = normalizeOrder(in.Order)
	if ids := parseIDList(in.IDs); len(ids) > 0 {
		q.IDs = ids
	}
	q.Keywords = strings.TrimSpace(in.Keywords)
	return q
}

func (q moduleQuery) frontInput(catNames, tagNames []string, size, kind int, shuffle bool) service.FrontListInput {
	in := service.FrontListInput{
		Keyword: q.Keywords, Categories: catNames, Tags: tagNames,
		Kind: kind, Sort: 1, Page: 1, Size: size, Ids: q.IDs,
		Shuffle: shuffle || q.Order == "rand",
	}
	if q.Order == "hot" {
		in.Sort = 0
	}
	return in
}
