package kit

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"net/url"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const totpIssuer = "子后台"

// GenerateTOTP 生成新的谷歌验证器密钥与扫码图。
func GenerateTOTP(account string) (secret, otpauthURI, qrDataURI string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      totpIssuer,
		AccountName: account,
		Period:      30,
		SecretSize:  20,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", "", "", err
	}
	qr, err := totpQRDataURI(key)
	if err != nil {
		return "", "", "", err
	}
	return key.Secret(), key.URL(), qr, nil
}

// TOTPQRFromSecret 用已有密钥重新生成 otpauth 与扫码图(待绑定重试时复用同一密钥)。
func TOTPQRFromSecret(account, secret string) (otpauthURI, qrDataURI string, err error) {
	if account == "" || secret == "" {
		return "", "", fmt.Errorf("totp account/secret empty")
	}
	uri := fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=6&period=30",
		url.PathEscape(totpIssuer), url.PathEscape(account), secret, url.QueryEscape(totpIssuer),
	)
	key, err := otp.NewKeyFromURL(uri)
	if err != nil {
		return "", "", err
	}
	qr, err := totpQRDataURI(key)
	if err != nil {
		return "", "", err
	}
	return key.URL(), qr, nil
}

// ValidateTOTP 校验 6 位谷歌验证器动态码。
func ValidateTOTP(code, secret string) bool {
	if code == "" || secret == "" {
		return false
	}
	ok, err := totp.ValidateCustom(code, secret, totpNow(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return err == nil && ok
}

func totpQRDataURI(key *otp.Key) (string, error) {
	img, err := key.Image(200, 200)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
