package ilua

import (
	"reflect"
	"testing"
	"time"

	"github.com/iglev/ilua/log"
	glua "github.com/yuin/gopher-lua"
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

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func TestLuar(t *testing.T) {
	if true {
		return
	}

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
	L.RegMod("mymod", LStateMod{
		"func1":     ModFunc1,
		"Num":       1234,
		"StringVal": "stringtest",
		"fib":       fib,
	})

	err := L.DoProFiles("./script/args.lua")
	if err != nil {
		log.Error("err=%v", err)
		return
	}

	now3 := time.Now()
	res3 := fib(35)
	log.Info("res=%v err=%v cost=%v", res3, "", time.Since(now3))

	now := time.Now()
	res, resErr := L.Call("mymod.fib", 35)
	log.Info("res=%v err=%v cost=%v", res, resErr, time.Since(now))

	// now2 := time.Now()
	// res2, resErr2 := L.Call("SVC_mod1.fib", 35)
	// log.Info("res2=%v err2=%v cost2=%v", res2, resErr2, time.Since(now2))

	// htFunc(L, func() {
	// 	// ...
	// })
}

type ltb struct {
	Name     string `ltb:"name"`
	Degree   bool   `ltb:"degree"`
	LogLevel int    `ltb:"log_level"`
	Sub      subltb `ltb:"sub"`
	Va       int
}

type subltb struct {
	N string `ltb:"n"`
}

func TestLTB(t *testing.T) {
	L := NewState()
	defer L.Close()
	tb := ltb{}
	ltbtype := reflect.TypeOf(tb)
	size := ltbtype.NumField()
	for i := 0; i < size; i++ {
		field := ltbtype.Field(i)
		log.Info("name=%v", field.Name)
		v, ok := field.Tag.Lookup("json")
		log.Info("v=%v ok=%v", v, ok)
		log.Info("tag=%v", field.Tag.Get("json"))
	}

	err := L.L().DoFile("./script/conf.lua")
	if err != nil {
		log.Error("err=%v", err)
		return
	}
	val, ok := L.L().Get(-1).(*glua.LTable)
	if !ok {
		log.Error("not ok")
		return
	}
	L.L().Pop(1)
	log.Info("dofile return-------------------%v", val)
}
