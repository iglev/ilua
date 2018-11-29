package ilua

import (
	"fmt"
	"testing"
)

// TestRunning ...
func TestRunning(t *testing.T) {
	L := NewState()
	defer Close(L)
	err := L.DoFile("test.lua")
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}
	fmt.Printf("test success\n")
}

func TestLogger(t *testing.T) {
	logerror("err=%v", 1234)
	loginfo("info=%v", "abcd")
}
