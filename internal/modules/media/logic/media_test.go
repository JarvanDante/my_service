package logic

import "testing"

func TestDetectContentType(t *testing.T) {
	cases := []struct {
		name, header, want string
	}{
		{"a.jpg", "", "image/jpeg"},
		{"a.MP4", "application/octet-stream", "video/mp4"},
		{"a.bin", "video/mp4", "video/mp4"},
	}
	for _, c := range cases {
		got := detectContentType(c.name, c.header)
		if got != c.want {
			t.Fatalf("%s header=%q got=%s want=%s", c.name, c.header, got, c.want)
		}
	}
}

func TestBuildObjectKey(t *testing.T) {
	key := buildObjectKey("video", "demo.mp4")
	if key == "" || len(key) < 10 {
		t.Fatalf("unexpected key: %s", key)
	}
	if key[len(key)-4:] != ".mp4" {
		t.Fatalf("ext missing: %s", key)
	}
}
