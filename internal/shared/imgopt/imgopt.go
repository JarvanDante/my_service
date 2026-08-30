// Package imgopt 统一存储图片在加密前按用途缩图/转 JPEG，避免原图直入。
package imgopt

import (
	"bytes"
	"image"
	"strings"

	"github.com/disintegration/imaging"
)

const quality = 80
const skipBytes = 300 * 1024

type spec struct {
	maxW int
	maxH int
}

func specFor(purpose string) (spec, bool) {
	switch strings.ToLower(strings.TrimSpace(purpose)) {
	case "avatar":
		return spec{512, 512}, true
	case "cover", "ad", "post", "image":
		return spec{1600, 1600}, true
	default:
		return spec{}, false
	}
}

// Compress 按用途缩到限制内并转 JPEG。GIF / 已足够小的 JPEG 原样返回。
func Compress(raw []byte, purpose string) (out []byte, changed bool) {
	s, ok := specFor(purpose)
	if !ok || len(raw) == 0 {
		return raw, false
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		return raw, false
	}
	if format == "gif" {
		return raw, false
	}
	needResize := cfg.Width > s.maxW || cfg.Height > s.maxH
	smallJPEG := !needResize && len(raw) <= skipBytes && (format == "jpeg" || format == "jpg")
	if smallJPEG {
		return raw, false
	}

	img, err := imaging.Decode(bytes.NewReader(raw), imaging.AutoOrientation(true))
	if err != nil {
		return raw, false
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w > s.maxW || h > s.maxH {
		img = imaging.Fit(img, s.maxW, s.maxH, imaging.Lanczos)
	}

	var buf bytes.Buffer
	if err = imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
		return raw, false
	}
	enc := buf.Bytes()
	if !needResize && len(enc) >= len(raw) {
		return raw, false
	}
	return enc, true
}

func JpegName(filename string) string {
	name := strings.TrimSpace(filename)
	if name == "" {
		return "image.jpg"
	}
	if i := strings.LastIndex(name, "."); i > 0 {
		return name[:i] + ".jpg"
	}
	return name + ".jpg"
}
