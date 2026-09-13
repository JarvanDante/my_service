package aesapi

import (
	"bytes"
	"testing"
)

func TestBase64RoundTrip(t *testing.T) {
	plain := []byte(`{"code":0,"message":"ok","data":{"id":1}}`)
	enc, err := EncryptBase64(plain)
	if err != nil {
		t.Fatal(err)
	}
	if enc == string(plain) {
		t.Fatal("ciphertext should differ")
	}
	got, err := DecryptBase64(enc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("got %s want %s", got, plain)
	}
}

func TestKnownVector(t *testing.T) {
	defer SetKey(DefaultKey)
	SetKey(DefaultKey)
	enc, err := EncryptBase64([]byte("hello-api"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecryptBase64(enc)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello-api" {
		t.Fatalf("got %q", got)
	}
	if enc != "bDaEQQA7U0vsgogEDJxyKw==" {
		t.Fatalf("known vector got %q", enc)
	}
}

func TestLooksLikeJSON(t *testing.T) {
	if !LooksLikeJSON([]byte("  {\"a\":1}")) || !LooksLikeJSON([]byte("[1]")) {
		t.Fatal("json")
	}
	if LooksLikeJSON([]byte("G0n0pQ==")) {
		t.Fatal("base64 is not json")
	}
}

func TestSetKeyFallback(t *testing.T) {
	defer SetKey(DefaultKey)
	SetKey("short")
	if ActiveKey() != DefaultKey {
		t.Fatal(ActiveKey())
	}
}
