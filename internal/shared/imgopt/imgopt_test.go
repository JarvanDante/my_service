package imgopt

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"
)

func jpegBytes(t *testing.T, w, h int, q int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 80, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: q}); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func pngBytes(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestCompressSkipSmallJPEG(t *testing.T) {
	raw := jpegBytes(t, 80, 80, 80)
	out, changed := Compress(raw, "cover")
	if changed {
		t.Fatal("small jpeg should stay")
	}
	if !bytes.Equal(out, raw) {
		t.Fatal("bytes changed")
	}
}

func TestCompressResizeCover(t *testing.T) {
	raw := jpegBytes(t, 2400, 1200, 95)
	out, changed := Compress(raw, "cover")
	if !changed {
		t.Fatal("oversized cover should compress")
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(out))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width > 1600 || cfg.Height > 1600 {
		t.Fatalf("still too big: %dx%d", cfg.Width, cfg.Height)
	}
}

func TestCompressAvatarBox(t *testing.T) {
	raw := jpegBytes(t, 1200, 800, 90)
	out, changed := Compress(raw, "avatar")
	if !changed {
		t.Fatal("avatar should fit 512")
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(out))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width > 512 || cfg.Height > 512 {
		t.Fatalf("avatar %dx%d", cfg.Width, cfg.Height)
	}
}

func TestCompressSkipGIF(t *testing.T) {
	img := image.NewPaletted(image.Rect(0, 0, 20, 20), color.Palette{color.White, color.Black})
	var buf bytes.Buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		t.Fatal(err)
	}
	raw := buf.Bytes()
	out, changed := Compress(raw, "post")
	if changed || !bytes.Equal(out, raw) {
		t.Fatal("gif should stay")
	}
}

func TestCompressSkipVideoPurpose(t *testing.T) {
	raw := jpegBytes(t, 2400, 1200, 90)
	out, changed := Compress(raw, "video")
	if changed || !bytes.Equal(out, raw) {
		t.Fatal("video purpose should not compress")
	}
}

func TestJpegName(t *testing.T) {
	if JpegName("a.PNG") != "a.jpg" {
		t.Fatal(JpegName("a.PNG"))
	}
	if JpegName("") != "image.jpg" {
		t.Fatal(JpegName(""))
	}
}
