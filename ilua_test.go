package ilua

import (
	"fmt"
	"testing"
	"time"

	glua "github.com/yuin/gopher-lua"
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

func TestLoadLibs(t *testing.T) {
	for i := 0; i < 1; i++ {
		func() {
			loginfo("-----------------------------------")
			L := NewState()
			defer Close(L)
			err := LoadLibs(L, "./test_script/args.lua")
			if err != nil {
				logerror("err=%v", err)
				return
			}
			loginfo("LoadLibs success")
			curr := 0
			timer := time.NewTimer(1 * time.Second)
		Loop:
			for {
				select {
				case <-timer.C:
					curr++
					err = L.CallByParam(glua.P{
						Fn:      L.GetGlobal("HotFix"),
						NRet:    1,
						Protect: true,
					})
					timer.Reset(1 * time.Second)
				default:
					if curr > 100 {
						timer.Stop()
						break Loop
					}
				}
			}

			loginfo("-----------------------------------\n")
		}()
	}
}
