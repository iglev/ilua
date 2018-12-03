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

type Node struct {
	Num int
}

type Person struct {
	Name string
	Age  int64
}

func (p *Person) Print() {
	log.Info("call Print, %v", p)
}

func (p *Person) PrintNode(no *Node) {
	log.Info("no=%v", no)
}

func (p *Person) GenNode(num int) *Node {
	return &Node{Num: num}
}

func ModFunc1() {
	log.Info("ModFunc1")
}

func TestLuar(t *testing.T) {
	L := NewState()
	defer L.Close()

	/*
		-- lua code, test RegType
		p = Person()
		p.Name = "testname"
		p.Age= 10
		p:Print()

		no = Node()
		no.Num = 111
		p:PrintNode(no) -- use custom type 'Node' as param
		p:PrintNode(p:GenNode(1234)) -- return custom type 'Node' to func 'PrintNode'
	*/
	L.RegType("Person", Person{})
	L.RegType("Node", Node{})

	/*
		-- lua code, test RegMod
		mymod.func1()
		LogInfo("mymod.Num=%v mymod.StringVal=%v", mymod.Num, mymod.StringVal)
	*/
	L.RegMod("mymod", map[string]interface{}{
		"func1":     ModFunc1,
		"Num":       1234,
		"StringVal": "stringtest",
	})

	err := L.DoFiles("./script/args.lua")
	if err != nil {
		log.Error("err=%v", err)
		return
	}
	htFunc(L, func() {
		// ...
	})
}
