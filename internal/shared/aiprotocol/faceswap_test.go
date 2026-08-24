package aiprotocol

import "testing"

func TestParseObjectRefURL(t *testing.T) {
	ref, ok := ParseObjectRef("http://127.0.0.1:19000/my-storage/my/image/abc/face.bnc?x=1", "")
	if !ok {
		t.Fatal("parse")
	}
	if ref.Bucket != "my-storage" || ref.Key != "my/image/abc/face.bnc" {
		t.Fatalf("%+v", ref)
	}
}

func TestParseObjectRefBareKey(t *testing.T) {
	ref, ok := ParseObjectRef("my/image/abc/face.bnc", "my-storage")
	if !ok || ref.Bucket != "my-storage" || ref.Key != "my/image/abc/face.bnc" {
		t.Fatalf("%v %+v", ok, ref)
	}
}

func TestObjectRefFromMap(t *testing.T) {
	ref, ok := ObjectRefFromAny(map[string]any{"bucket": "my-storage", "key": "a/b.bnc"}, "")
	if !ok || !ref.OK() {
		t.Fatal(ref)
	}
}
