package logic

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/JarvanDante/my_service/internal/modules/comics/service"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

// moduleQuery 模块检索条件。对齐公司 JSON，底层查 Postgres comics 表。
type moduleQuery struct {
	TagIDs    []int64
	CatIDs    []int64
	Order     string
	IsEnd     string
	PayType   int
	IDs       []int64
	Keywords  string
	Recommend bool
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

func parseFlag(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	s := strings.TrimSpace(strings.Trim(string(raw), `"`))
	s = strings.ToLower(s)
	switch s {
	case "y", "yes", "true", "1":
		return "y"
	case "n", "no", "false", "0":
		return "n"
	default:
		return ""
	}
}

func parsePayType(raw json.RawMessage) int {
	if len(raw) == 0 {
		return 0
	}
	s := strings.ToLower(strings.TrimSpace(strings.Trim(string(raw), `"`)))
	switch s {
	case "1", "vip":
		return 1
	case "2", "coin", "paid", "money":
		return 2
	case "3", "free":
		return 3
	default:
		return 0
	}
}

func normalizeOrder(order string) string {
	switch strings.ToLower(strings.TrimSpace(order)) {
	case "rand", "random":
		return "rand"
	case "hot", "click", "ranking":
		return "hot"
	case "like":
		return "like"
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
	m := map[string]any{}
	if s := encodeIDCSV(q.TagIDs); s != "" {
		m["tag_id"] = s
	}
	if s := encodeIDCSV(q.CatIDs); s != "" {
		m["cat_id"] = s
	}
	order := normalizeOrder(q.Order)
	m["order"] = order
	if q.IsEnd == "y" || q.IsEnd == "n" {
		m["is_end"] = q.IsEnd
	}
	switch q.PayType {
	case 1:
		m["pay_type"] = "vip"
	case 2:
		m["pay_type"] = "coin"
	case 3:
		m["pay_type"] = "free"
	}
	if s := encodeIDCSV(q.IDs); s != "" {
		m["ids"] = s
	}
	if kw := strings.TrimSpace(q.Keywords); kw != "" {
		m["keywords"] = kw
	}
	if q.Recommend {
		m["recommend"] = "y"
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
		TagID     json.RawMessage `json:"tag_id"`
		CatID     json.RawMessage `json:"cat_id"`
		Order     string          `json:"order"`
		IsEnd     json.RawMessage `json:"is_end"`
		PayType   json.RawMessage `json:"pay_type"`
		IDs       json.RawMessage `json:"ids"`
		Keywords  string          `json:"keywords"`
		Recommend json.RawMessage `json:"recommend"`
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
	if flag := parseFlag(in.IsEnd); flag != "" {
		q.IsEnd = flag
	}
	q.PayType = parsePayType(in.PayType)
	if ids := parseIDList(in.IDs); len(ids) > 0 {
		q.IDs = ids
	}
	q.Keywords = strings.TrimSpace(in.Keywords)
	q.Recommend = parseFlag(in.Recommend) == "y"
	return q
}

func (q moduleQuery) listFilter(catNames, tagNames []string, size int, shuffle bool) service.ListFilter {
	f := service.ListFilter{
		Categories:    catNames,
		Tags:          tagNames,
		Keyword:       q.Keywords,
		PayType:       q.PayType,
		OnlyRecommend: q.Recommend,
		Ids:           q.IDs,
		Page:          1,
		Size:          size,
		Sort:          2,
		Shuffle:       shuffle || q.Order == "rand",
	}
	switch q.Order {
	case "hot":
		f.Sort = 1
	case "like":
		f.Sort = 3
	}
	if q.IsEnd == "y" {
		f.UpdateStatus = 2
	} else if q.IsEnd == "n" {
		f.UpdateStatus = 1
	}
	return f
}
