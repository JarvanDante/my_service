package logic

import (
	"testing"

	"github.com/JarvanDante/my_service/internal/modules/aitask/service"
)

func TestSetsOf(t *testing.T) {
	t.Parallel()
	if got := setsOf(nil); got != 1 {
		t.Fatalf("nil params = %d", got)
	}
	if got := setsOf(map[string]any{"sets": float64(3)}); got != 3 {
		t.Fatalf("float sets = %d", got)
	}
	if got := setsOf(map[string]any{"sets": 0}); got != 1 {
		t.Fatalf("zero sets should fallback, got %d", got)
	}
}

func TestTaskNeedUserJoin(t *testing.T) {
	t.Parallel()
	if taskNeedUserJoin(service.TaskFilter{}) {
		t.Fatal("empty filter should not join users")
	}
	if !taskNeedUserJoin(service.TaskFilter{Nickname: "张"}) {
		t.Fatal("nickname filter should join users")
	}
	if !taskNeedUserJoin(service.TaskFilter{DeviceType: "h5"}) {
		t.Fatal("device filter should join users")
	}
}
