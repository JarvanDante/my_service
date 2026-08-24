package aiprotocol

import (
	"net/url"
	"strings"
)

// ParseObjectRef 把公开 URL 或裸 key 收成 bucket/key。
func ParseObjectRef(raw, defaultBucket string) (ObjectRef, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ObjectRef{}, false
	}
	if defaultBucket == "" {
		defaultBucket = DefaultBucket
	}
	if i := strings.IndexAny(raw, "?#"); i >= 0 {
		raw = raw[:i]
	}
	if !strings.Contains(raw, "://") {
		return ObjectRef{Bucket: defaultBucket, Key: strings.TrimLeft(raw, "/")}, true
	}
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" || u.Path == "/" {
		return ObjectRef{}, false
	}
	path := strings.Trim(u.Path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return ObjectRef{Bucket: parts[0], Key: parts[1]}, true
	}
	if path != "" {
		return ObjectRef{Bucket: defaultBucket, Key: path}, true
	}
	return ObjectRef{}, false
}

// ObjectRefFromAny 认 {bucket,key} 对象或 URL/key 字符串。
func ObjectRefFromAny(v any, defaultBucket string) (ObjectRef, bool) {
	switch t := v.(type) {
	case ObjectRef:
		if t.Key == "" {
			return ObjectRef{}, false
		}
		if t.Bucket == "" {
			if defaultBucket == "" {
				defaultBucket = DefaultBucket
			}
			t.Bucket = defaultBucket
		}
		return t, true
	case map[string]any:
		bucket, _ := t["bucket"].(string)
		key, _ := t["key"].(string)
		bucket = strings.TrimSpace(bucket)
		key = strings.TrimSpace(key)
		if key == "" {
			return ObjectRef{}, false
		}
		if bucket == "" {
			if defaultBucket == "" {
				defaultBucket = DefaultBucket
			}
			bucket = defaultBucket
		}
		return ObjectRef{Bucket: bucket, Key: key}, true
	case string:
		return ParseObjectRef(t, defaultBucket)
	default:
		return ObjectRef{}, false
	}
}

func FirstString(params map[string]any, keys ...string) string {
	if params == nil {
		return ""
	}
	for _, k := range keys {
		if v, ok := params[k]; ok {
			if s, ok := v.(string); ok {
				if s = strings.TrimSpace(s); s != "" {
					return s
				}
			}
		}
	}
	return ""
}
