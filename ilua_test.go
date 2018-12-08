package ilua

import (
	"testing"
	"time"

	"github.com/iglev/ilua/log"
	glua "github.com/yuin/gopher-lua"
)

type Person struct {
	Name string
	Age  int
	BD   BirthDay
}

type BirthDay struct {
	Y int16
	M uint8
	D uint8
}

func TestRegType(t *testing.T) {
	L := NewState()
	defer L.Close()
	L.RegType("Person", Person{})
	err := L.L().DoString(`
		p = Person()
		p.Name = "testname"
		p.Age = 12
		p.BD = { Y=2000, M=1, D=1 }
		print(p.Name, p.Age, p.BD.Y, p.BD.M, p.BD.D)
		print(type(p.BD)) -- userdata
	`)
	if err != nil {
		log.Error("err=%v", err)
		return
	}
}

func createPerson(name string, age int) *Person {
	return &Person{Name: name, Age: age}
}

func TestCallFunc(t *testing.T) {
	L := NewState()
	defer L.Close()
	doErr := L.DoProFiles("./script/args.lua")
	if doErr != nil {
		log.Error("err=%v", doErr)
		return
	}
	L.RegMod("mymod", LStateMod{
		"create": createPerson,
		"incNum": 10,
	})
	ret, err := L.Call("mymod.create", "tname", 12)
	if err != nil {
		log.Error("err=%v", err)
		return
	}
	res, resOK := ret.(*glua.LUserData)
	if !resOK {
		log.Error("not userdata, ret=%v", ret)
		return
	}
	log.Info("res=%v", res.Value)
}

type ltb struct {
	Name     string
	Degree   bool
	LogLevel int
	Sub      subltb
	Va       int
}

type subltb struct {
	N  string
	XY subltb2
}

type subltb2 struct {
	X, Y int64
}

func TestLTB(t *testing.T) {
	L := NewState()
	defer L.Close()
	// res, _ := L.UnmarshalLTB("./script/conf.lua", ltb{})
	res, _ := UnmarshalLTB("./script/conf.lua", ltb{})
	log.Info("res=%+v", res)
}

func TestHotfix(t *testing.T) {
	L := NewState(SetHotfix(true, false))
	defer L.Close()
	doErr := L.DoProFiles("./script/args.lua")
	if doErr != nil {
		log.Error("err=%v", doErr)
		return
	}
	count := 0
	L.RegMod("mymod", LStateMod{
		"func1": func() {
			count++
			// log.Info("call mymod.func1")
		},
	})
	/*
		timer := time.NewTimer(1 * time.Second)
		for {
			select {
			case <-timer.C:
				func() {
					_, err := L.Call("mymod.func1")
					if err != nil {
						log.Error("err=%v", err)
						return
					}
				}()
				timer.Reset(1 * time.Second)
			}
		}
	*/
	curr := time.Now()
	for i := 0; i < 1000000; i++ {
		_, err := L.Call("mymod.func1")
		if err != nil {
			log.Error("err=%v", err)
			return
		}
	}
	log.Info("cost=%v", time.Since(curr))
}
