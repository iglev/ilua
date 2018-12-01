package ilua

import (
	"testing"
	"time"

	"github.com/iglev/ilua/log"
)

// // TestRunning ...
// func TestRunning(t *testing.T) {
// 	L := NewState()
// 	defer L.Close()
// 	err := L.L().DoFile("test.lua")
// 	if err != nil {
// 		fmt.Printf("err=%v\n", err)
// 		return
// 	}
// 	fmt.Printf("test success\n")
// }

// func TestLogger(t *testing.T) {
// 	log.Error("err=%v", 1234)
// 	log.Info("info=%v", "abcd")
// }

func TestLoadLibs(t *testing.T) {
	L := NewState()
	defer L.Close()
	err := LoadLibs(L, "./script/args.lua")
	if err != nil {
		log.Error("err=%v", err)
		return
	}
	timer := time.NewTimer(1 * time.Second)
	c := 1000
Loop:
	for {
		select {
		case <-timer.C:
			if c < 0 {
				timer.Stop()
				break Loop
			}
			c--
			err = L.CheckHotfix()
			if err != nil {
				log.Error("err=%v", err)
				return
			}
			timer.Reset(1 * time.Second)
		}
	}
	log.Info("LoadLibs success")
}
