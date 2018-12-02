package ilua

import (
	"testing"
	"time"

	"github.com/iglev/ilua/log"
)

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

func luafn(str string, other ...interface{}) string {
	log.Info("str=%v ohter=%v", str, other)
	return "abcd"
}

type Person struct {
	Name string
	Age  int64
}

func (p *Person) Print() {
	log.Info("call Print, %v", p)
}

func TestLuar(t *testing.T) {
	L := NewState()
	defer L.Close()

	/*
		-- lua code
		p = Person()
		p.Name = "testname"
		p.Age= 10
		p:Print()
	*/
	L.RegType("Person", Person{})

	err := L.DoFiles("./script/args.lua")
	if err != nil {
		log.Error("err=%v", err)
		return
	}
	htFunc(L, func() {
		// ...
	})
}
