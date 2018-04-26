package contract

import (
	"testing"
	"time"

	".."
)

func TestMock_IsAllowedAt(t *testing.T) {
	m := Mock{}
	ok, _ := m.IsAllowedAt(tür.MieterID("franz"), time.Now())
	if !ok {
		t.Error("Sollte im mock immer erlaubt sein")
	}
}
