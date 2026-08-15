package kit

import (
	"strings"
)

// userIdAlphabet 与 dm-php UserService::ENCODE_STR 一致。
// 前台「编号」是数字 id 的短码, 不是数据库主键。
const userIdAlphabet = "SCUWDG3HE859QA4B1NOPIV67XLYFJ2RZTKM"

// EncodeUserId 通过 userId 生成公开编号/邀请码, 对齐 dm-php encodeUserId。
func EncodeUserId(userId int64) string {
	if userId <= 0 {
		return ""
	}
	sLen := int64(len(userIdAlphabet))
	num := userId
	var code []byte
	for num > 0 {
		mod := num % sLen
		num = (num - mod) / sLen
		code = append([]byte{userIdAlphabet[mod]}, code...)
	}
	if len(code) < 4 {
		padded := []byte("0000")
		copy(padded[4-len(code):], code)
		return string(padded)
	}
	return string(code)
}

// DecodeUserId 通过公开编号还原 userId, 对齐 dm-php decodeUserId。
func DecodeUserId(code string) int64 {
	code = strings.TrimSpace(code)
	if code == "" {
		return 0
	}
	if i := strings.LastIndex(code, "0"); i >= 0 {
		code = code[i+1:]
	}
	if code == "" {
		return 0
	}
	sLen := int64(len(userIdAlphabet))
	runes := []byte(code)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	var num int64
	pow := int64(1)
	for _, ch := range runes {
		pos := strings.IndexByte(userIdAlphabet, ch)
		if pos < 0 {
			return 0
		}
		num += int64(pos) * pow
		pow *= sLen
	}
	return num
}
