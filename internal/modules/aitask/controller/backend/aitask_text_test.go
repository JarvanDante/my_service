package backend

import "testing"

func TestBizTypeText(t *testing.T) {
	t.Parallel()
	if got := bizTypeText(1, nil); got != "图片换脸" {
		t.Fatalf("photo = %s", got)
	}
	if got := bizTypeText(1, map[string]any{"media_type": "video"}); got != "视频换脸" {
		t.Fatalf("video = %s", got)
	}
	if got := bizTypeText(4, nil); got != "图生视频" {
		t.Fatalf("i2v = %s", got)
	}
}

func TestStatusText(t *testing.T) {
	t.Parallel()
	if got := statusText(1); got != "待处理" {
		t.Fatalf("queued = %s", got)
	}
	if got := statusText(3); got != "成功" {
		t.Fatalf("ok = %s", got)
	}
	if got := statusText(5); got != "退款" {
		t.Fatalf("refund = %s", got)
	}
}
