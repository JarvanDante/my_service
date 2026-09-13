package middleware

import "testing"

func TestIsFrontAPI(t *testing.T) {
	if !isFrontAPI("/front/v1/user/login") || isFrontAPI("/backend/admin/login") || isFrontAPI("/health") {
		t.Fatal("front prefix")
	}
}

func TestIsCryptoSkipPath(t *testing.T) {
	if !isCryptoSkipPath("/front/v1/media/upload") ||
		!isCryptoSkipPath("/front/v1/media/object") ||
		!isCryptoSkipPath("/front/v1/media/multipart/part") {
		t.Fatal("skip")
	}
	if isCryptoSkipPath("/front/v1/media/multipart/init") || isCryptoSkipPath("/front/v1/user/login") {
		t.Fatal("json apis must encrypt")
	}
}
