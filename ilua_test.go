package ilua

import (
	"testing"
	"time"

	"github.com/iglev/ilua/log"
	"github.com/iglev/ilua/luar"
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

// func TestLoadLibs(t *testing.T) {
// 	L := NewState()
// 	defer L.Close()
// 	err := LoadLibs(L, "./script/args.lua")
// 	if err != nil {
// 		log.Error("err=%v", err)
// 		return
// 	}
// 	timer := time.NewTimer(1 * time.Second)
// 	c := 1000
// Loop:
// 	for {
// 		select {
// 		case <-timer.C:
// 			if c < 0 {
// 				timer.Stop()
// 				break Loop
// 			}
// 			c--
// 			err = L.CheckHotfix()
// 			if err != nil {
// 				log.Error("err=%v", err)
// 				return
// 			}
// 			timer.Reset(1 * time.Second)
// 		}
// 	}
// 	log.Info("LoadLibs success")
// }

func htFunc(L *LState, cb func()) {
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
			err := L.CheckHotfix()
			if err != nil {
				log.Error("err=%v", err)
				return
			}
			cb()
			timer.Reset(1 * time.Second)
		}
	}
}

type Node struct {
	Name  string
	token string
	Num   int
}

func (no *Node) SetToken(t string) {
	no.token = t
}

func (no *Node) Token() string {
	return no.token
}

func TestLuar(t *testing.T) {
	L := NewState()
	defer L.Close()

	no := &Node{
		Name:  "myname",
		token: "token",
		Num:   123,
	}
	tmp := luar.New(L.L(), no)
	log.Info("tmp=%v", tmp)
	L.L().SetGlobal("no", tmp)

	err := LoadLibs(L, "./script/args.lua")
	if err != nil {
		log.Error("err=%v", err)
		return
	}
	htFunc(L, func() {
		log.Info("golag no=%v", no)
	})
}
