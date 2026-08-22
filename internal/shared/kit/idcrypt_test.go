package kit

import "testing"

func TestEncodePublicId(t *testing.T) {
	cases := map[int64]string{
		1:     "XENZBw",
		2:     "X0NZTQ",
		3:     "XkNaRg",
		6:     "W0NdUQ",
		10:    "XElSDQI",
		35:    "XkxSQBU",
		100:   "XElYD0RUCQ",
		1234:  "XEtbAUxFAhk",
		99999: "VEBRDE8LWhgfEVI",
	}
	for id, want := range cases {
		if got := EncodePublicId(id); got != want {
			t.Fatalf("EncodePublicId(%d)=%s want %s", id, got, want)
		}
	}
	if EncodePublicId(0) != "" {
		t.Fatal("id=0 should be empty")
	}
}

func TestDecodePublicIdRoundTrip(t *testing.T) {
	ids := []int64{1, 2, 3, 6, 10, 35, 100, 1234, 99999}
	for _, id := range ids {
		code := EncodePublicId(id)
		if got := DecodePublicId(code); got != id {
			t.Fatalf("DecodePublicId(%s)=%d want %d", code, got, id)
		}
	}
}
