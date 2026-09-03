package logic

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"strings"
)

// 去掉 0/O/1/I/L，降低误认。
const charset = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"

const (
	imgW   = 120
	imgH   = 40
	scale  = 3
	glyphW = 5
	glyphH = 7
)

// 5x7，每行 5bit，最高位为左。
var glyphs = map[byte][7]byte{
	'2': {0x0E, 0x11, 0x01, 0x02, 0x04, 0x08, 0x1F},
	'3': {0x0E, 0x11, 0x01, 0x06, 0x01, 0x11, 0x0E},
	'4': {0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02},
	'5': {0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E},
	'6': {0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E},
	'7': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08},
	'8': {0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E},
	'9': {0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C},
	'A': {0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'B': {0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E},
	'C': {0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E},
	'D': {0x1E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x1E},
	'E': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F},
	'F': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x10},
	'G': {0x0E, 0x11, 0x10, 0x13, 0x11, 0x11, 0x0E},
	'H': {0x11, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'J': {0x01, 0x01, 0x01, 0x01, 0x11, 0x11, 0x0E},
	'K': {0x11, 0x12, 0x14, 0x18, 0x14, 0x12, 0x11},
	'M': {0x11, 0x1B, 0x15, 0x15, 0x11, 0x11, 0x11},
	'N': {0x11, 0x19, 0x15, 0x13, 0x11, 0x11, 0x11},
	'P': {0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10},
	'Q': {0x0E, 0x11, 0x11, 0x11, 0x15, 0x12, 0x0D},
	'R': {0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11},
	'S': {0x0E, 0x11, 0x10, 0x0E, 0x01, 0x11, 0x0E},
	'T': {0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04},
	'U': {0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
	'V': {0x11, 0x11, 0x11, 0x11, 0x11, 0x0A, 0x04},
	'W': {0x11, 0x11, 0x11, 0x15, 0x15, 0x1B, 0x11},
	'X': {0x11, 0x11, 0x0A, 0x04, 0x0A, 0x11, 0x11},
	'Y': {0x11, 0x11, 0x0A, 0x04, 0x04, 0x04, 0x04},
	'Z': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x10, 0x1F},
}

func normalizeCode(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func codesEqual(got, want string) bool {
	return normalizeCode(got) == normalizeCode(want)
}

func pickCode() string {
	raw := make([]byte, 4)
	if _, err := rand.Read(raw); err != nil {
		return "A2B3"
	}
	out := make([]byte, 4)
	for i := range out {
		out[i] = charset[int(raw[i])%len(charset)]
	}
	return string(out)
}

func rndByte() byte {
	var b [1]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	return b[0]
}

func rndN(n int) int {
	if n <= 0 {
		return 0
	}
	return int(rndByte()) % n
}

func renderDataURI(code string) (string, error) {
	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	bg := color.RGBA{R: 22, G: 18, B: 32, A: 255}
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			img.Set(x, y, bg)
		}
	}

	gap := 8
	block := glyphW * scale
	total := 4*block + 3*gap
	originX := (imgW - total) / 2
	originY := (imgH - glyphH*scale) / 2

	for i := 0; i < len(code) && i < 4; i++ {
		ch := code[i]
		fg := color.RGBA{
			R: 160 + rndByte()%80,
			G: 140 + rndByte()%90,
			B: 220,
			A: 255,
		}
		ox := originX + i*(block+gap) + rndN(3) - 1
		oy := originY + rndN(5) - 2
		drawGlyph(img, ch, ox, oy, fg)
	}

	for i := 0; i < 5; i++ {
		c := color.RGBA{R: 80 + rndByte()%80, G: 70 + rndByte()%70, B: 120 + rndByte()%80, A: 180}
		drawLine(img, rndN(imgW), rndN(imgH), rndN(imgW), rndN(imgH), c)
	}
	for i := 0; i < 90; i++ {
		c := color.RGBA{R: 90 + rndByte()%120, G: 80 + rndByte()%120, B: 140 + rndByte()%100, A: 160}
		img.Set(rndN(imgW), rndN(imgH), c)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func drawGlyph(img *image.RGBA, ch byte, ox, oy int, fg color.RGBA) {
	rows, ok := glyphs[ch]
	if !ok {
		return
	}
	for row := 0; row < glyphH; row++ {
		bits := rows[row]
		for col := 0; col < glyphW; col++ {
			if bits&(1<<uint(4-col)) == 0 {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					x := ox + col*scale + dx
					y := oy + row*scale + dy
					if image.Pt(x, y).In(img.Bounds()) {
						img.Set(x, y, fg)
					}
				}
			}
		}
	}
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	for {
		if image.Pt(x0, y0).In(img.Bounds()) {
			img.Set(x0, y0, c)
		}
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
