package kit

import (
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func TestGenerateAndValidateTOTP(t *testing.T) {
	now := time.Date(2026, 9, 7, 16, 30, 0, 0, time.UTC)
	totpNow = func() time.Time { return now }
	t.Cleanup(func() { totpNow = time.Now })

	secret, uri, qr, err := GenerateTOTP("MY-子后台", "admin")
	if err != nil {
		t.Fatal(err)
	}
	if secret == "" || !strings.Contains(uri, "MY-%E5%AD%90%E5%90%8E%E5%8F%B0") && !strings.Contains(uri, "MY-子后台") {
		t.Fatalf("issuer missing in uri=%q", uri)
	}
	if !strings.Contains(uri, "admin") || !strings.HasPrefix(qr, "data:image/png;base64,") {
		t.Fatalf("unexpected totp payload secret=%q uri=%q qrPrefix=%q", secret, uri, qr[:min(32, len(qr))])
	}

	code, err := totp.GenerateCode(secret, now)
	if err != nil {
		t.Fatal(err)
	}
	if !ValidateTOTP(code, secret) {
		t.Fatal("current code should pass")
	}
	if ValidateTOTP("000000", secret) {
		t.Fatal("wrong code should fail")
	}
	if ValidateTOTP("", secret) || ValidateTOTP(code, "") {
		t.Fatal("empty code/secret should fail")
	}
}

func TestTOTPQRFromSecretReusesSameSecret(t *testing.T) {
	secret, _, _, err := GenerateTOTP("JH-子后台", "yyleader")
	if err != nil {
		t.Fatal(err)
	}
	uri, qr, err := TOTPQRFromSecret("JH-子后台", "yyleader", secret)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(uri, secret) || !strings.HasPrefix(qr, "data:image/png;base64,") {
		t.Fatalf("qr from secret mismatch uri=%q", uri)
	}
	if !strings.Contains(uri, "JH") {
		t.Fatalf("issuer should keep site prefix uri=%q", uri)
	}
	if _, _, err = TOTPQRFromSecret("JH-子后台", "", secret); err == nil {
		t.Fatal("empty account should fail")
	}
}

func TestSiteTotpIssuer(t *testing.T) {
	t.Setenv("SITE_CODE", "my")
	if got := SiteTotpIssuer(); got != "MY-子后台" {
		t.Fatalf("SiteTotpIssuer=%q", got)
	}
	t.Setenv("SITE_CODE", "jh")
	if got := SiteTotpIssuer(); got != "JH-子后台" {
		t.Fatalf("SiteTotpIssuer=%q", got)
	}
}
