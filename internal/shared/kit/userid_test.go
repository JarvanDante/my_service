package kit

import "testing"

func TestEncodeUserId(t *testing.T) {
	cases := map[int64]string{
		1:      "000C",
		2:      "000U",
		10:     "0005",
		34:     "000M",
		35:     "00CS",
		36:     "00CC",
		100:    "00UR",
		1234:   "0CS8",
		99999:  "U96D",
		100000: "U96G",
		350000: "EGLS",
	}
	for id, want := range cases {
		if got := EncodeUserId(id); got != want {
			t.Fatalf("EncodeUserId(%d)=%s want %s", id, got, want)
		}
	}
	if EncodeUserId(0) != "" {
		t.Fatal("id=0 should be empty")
	}
}

func TestDecodeUserIdRoundTrip(t *testing.T) {
	ids := []int64{1, 2, 10, 34, 35, 36, 100, 1234, 99999, 100000, 350000}
	for _, id := range ids {
		code := EncodeUserId(id)
		if got := DecodeUserId(code); got != id {
			t.Fatalf("DecodeUserId(%s)=%d want %d", code, got, id)
		}
	}
}
