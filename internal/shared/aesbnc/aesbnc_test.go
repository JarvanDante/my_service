package aesbnc

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	plain := []byte("hello-cover")
	enc, err := Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(enc, plain) {
		t.Fatal("ciphertext should differ")
	}
	got, err := Decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("got %q want %q", got, plain)
	}
}

func TestJpegNotLooksLikeImageAfterEncrypt(t *testing.T) {
	plain := append([]byte{0xFF, 0xD8, 0xFF, 0xE0}, bytes.Repeat([]byte("x"), 32)...)
	enc, err := Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	if LooksLikeImage(enc) {
		t.Fatal("ciphertext must not look like jpeg")
	}
	if !LooksLikeImage(plain) {
		t.Fatal("plain jpeg magic")
	}
}

func TestKnownVector(t *testing.T) {
	enc, err := Encrypt([]byte("hello-cover"))
	if err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(enc); got != "0c849758042af1c9174d3b41c36c7f02" {
		t.Fatalf("got %s", got)
	}
}

func TestSetKeyFallback(t *testing.T) {
	defer SetKey(DefaultKey)
	SetKey("short")
	if ActiveKey() != DefaultKey {
		t.Fatal(ActiveKey())
	}
	SetKey("1234567890123456")
	if ActiveKey() != "1234567890123456" {
		t.Fatal(ActiveKey())
	}
	enc, err := Encrypt([]byte("hello-cover"))
	if err != nil {
		t.Fatal(err)
	}
	SetKey(DefaultKey)
	if _, err := Decrypt(enc); err == nil {
		t.Fatal("old ciphertext must fail with default key")
	}
}

func TestToBncKey(t *testing.T) {
	if ToBncKey("a/cover.jpg") != "a/cover.bnc" {
		t.Fatal(ToBncKey("a/cover.jpg"))
	}
	if ToBncKey("a/cover.bnc") != "a/cover.bnc" {
		t.Fatal(ToBncKey("a/cover.bnc"))
	}
	if ToBncKey("") != "image.bnc" {
		t.Fatal(ToBncKey(""))
	}
}

func TestIsEncryptedName(t *testing.T) {
	if !IsEncryptedName("a/cover.bnc") || !IsEncryptedName("my/image/xxx/bnc") || !IsEncryptedName("bnc") {
		t.Fatal("expected ciphertext names")
	}
	if IsEncryptedName("a/cover.jpg") {
		t.Fatal("jpg is plaintext")
	}
}
