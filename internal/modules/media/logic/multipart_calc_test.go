package logic

import (
	"math"
	"testing"
)

func TestPartCountCalc(t *testing.T) {
	partSize := int64(8 << 20)
	size := int64(20 << 20) // 20MiB
	n := int(math.Ceil(float64(size) / float64(partSize)))
	if n != 3 {
		t.Fatalf("got %d want 3", n)
	}
}
