package logic

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strings"
	"testing"
)

func TestNormalizeAndEqual(t *testing.T) {
	if normalizeCode(" a2b3 ") != "A2B3" {
		t.Fatalf("normalizeCode = %q", normalizeCode(" a2b3 "))
	}
	if !codesEqual("a2b3", "A2B3") {
		t.Fatal("codesEqual should ignore case")
	}
	if codesEqual("A2B3", "A2B4") {
		t.Fatal("codesEqual should reject mismatch")
	}
}

func TestPickCode(t *testing.T) {
	seen := map[byte]bool{}
	for i := 0; i < 80; i++ {
		code := pickCode()
		if len(code) != 4 {
			t.Fatalf("pickCode len = %d", len(code))
		}
		for j := 0; j < 4; j++ {
			ch := code[j]
			if !strings.ContainsRune(charset, rune(ch)) {
				t.Fatalf("unexpected char %q", ch)
			}
			seen[ch] = true
		}
	}
	if len(seen) < 8 {
		t.Fatalf("pickCode too little variety: %d", len(seen))
	}
}

func TestRenderDataURI(t *testing.T) {
	uri, err := renderDataURI("A2B3")
	if err != nil {
		t.Fatal(err)
	}
	const prefix = "data:image/png;base64,"
	if !strings.HasPrefix(uri, prefix) {
		t.Fatalf("prefix = %q", uri[:min(len(uri), 24)])
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(uri, prefix))
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != imgW || img.Bounds().Dy() != imgH {
		t.Fatalf("size = %v", img.Bounds())
	}
}

func TestGlyphsCoverCharset(t *testing.T) {
	for i := 0; i < len(charset); i++ {
		ch := charset[i]
		if _, ok := glyphs[ch]; !ok {
			t.Fatalf("missing glyph %q", ch)
		}
		if _, err := renderDataURI(string([]byte{ch, ch, ch, ch})); err != nil {
			t.Fatalf("render %q: %v", ch, err)
		}
	}
}
