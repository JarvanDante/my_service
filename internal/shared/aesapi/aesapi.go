// Package aesapi 前台 JSON 接口传输加密: AES-128-ECB + PKCS7, 密文再 Base64。
// 与封面 .bnc（aesbnc）密钥分开, 只用于 /front/v1 请求体和响应信封。
package aesapi

import (
	"crypto/aes"
	"encoding/base64"
	"errors"
	"strings"
)

const DefaultKey = "9f3a6c1e8b4d0275"
const DefaultDebugKey = "myh5dbg7k2p9q4x1"
const DefaultDebugHeader = "c7e4a19b3f6820d54e8a16c2b9f735d1"

var activeKey = DefaultKey

func SetKey(k string) {
	k = strings.TrimSpace(k)
	if len(k) != aes.BlockSize {
		activeKey = DefaultKey
		return
	}
	activeKey = k
}

func ActiveKey() string { return activeKey }

func Encrypt(plain []byte) ([]byte, error) {
	if plain == nil {
		plain = []byte{}
	}
	block, err := aes.NewCipher([]byte(activeKey))
	if err != nil {
		return nil, err
	}
	padded := pkcs7Pad(plain, aes.BlockSize)
	out := make([]byte, len(padded))
	for i := 0; i < len(padded); i += aes.BlockSize {
		block.Encrypt(out[i:i+aes.BlockSize], padded[i:i+aes.BlockSize])
	}
	return out, nil
}

func Decrypt(data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, errors.New("密文长度无效")
	}
	block, err := aes.NewCipher([]byte(activeKey))
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	for i := 0; i < len(data); i += aes.BlockSize {
		block.Decrypt(out[i:i+aes.BlockSize], data[i:i+aes.BlockSize])
	}
	return pkcs7Unpad(out, aes.BlockSize)
}

func EncryptBase64(plain []byte) (string, error) {
	raw, err := Encrypt(plain)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}

func DecryptBase64(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "'\"")
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return Decrypt(raw)
}

func LooksLikeJSON(b []byte) bool {
	b = bytesTrimSpace(b)
	return len(b) > 0 && (b[0] == '{' || b[0] == '[')
}

func bytesTrimSpace(b []byte) []byte {
	i, j := 0, len(b)
	for i < j && (b[i] == ' ' || b[i] == '\n' || b[i] == '\r' || b[i] == '\t') {
		i++
	}
	for j > i && (b[j-1] == ' ' || b[j-1] == '\n' || b[j-1] == '\r' || b[j-1] == '\t') {
		j--
	}
	return b[i:j]
}

func pkcs7Pad(b []byte, n int) []byte {
	pad := n - (len(b) % n)
	out := make([]byte, len(b)+pad)
	copy(out, b)
	for i := len(b); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

func pkcs7Unpad(b []byte, n int) ([]byte, error) {
	if len(b) == 0 || len(b)%n != 0 {
		return nil, errors.New("填充无效")
	}
	pad := int(b[len(b)-1])
	if pad == 0 || pad > n || pad > len(b) {
		return nil, errors.New("填充无效")
	}
	for i := len(b) - pad; i < len(b); i++ {
		if b[i] != byte(pad) {
			return nil, errors.New("填充无效")
		}
	}
	return b[:len(b)-pad], nil
}
