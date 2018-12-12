package ilua

import (
	// "encoding/json"
	"os"
	"testing"
	"time"

	"github.com/iglev/ilua/log"
	json "github.com/json-iterator/go"
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
	Name     string `json:"name"`
	Degree   bool   `json:"degree"`
	LogLevel int    `json:"log_level"`
	Sub      subltb `json:"sub"`
}

type subltb struct {
	N  string  `json:"n"`
	XY subltb2 `json:"xy"`
}

type subltb2 struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

func ltbload(L *LState) {
	// res, _ := L.UnmarshalLTB("./script/conf.lua", ltb{})
	// log.Info("res=%v", res)
	L.UnmarshalLTB("./script/conf.lua", ltb{})
}

func TestLTB(t *testing.T) {
	L := NewState()
	defer L.Close()
	// res, _ := L.UnmarshalLTB("./script/conf.lua", ltb{})
	// res, _ := UnmarshalLTB("./script/conf.lua", ltb{})
	// log.Info("res=%v", res)
	now := time.Now()
	for i := 0; i < 10000; i++ {
		ltbload(L)
	}
	log.Info("ltb--------------------------cost=%v", time.Since(now))
}

func jsonload() {
	var one ltb
	file, err := os.Open("./script/conf.json")
	if err != nil {
		log.Error("not found file err=%v", err)
		return
	}
	defer file.Close()
	jsonparser := json.NewDecoder(file)
	err = jsonparser.Decode(&one)
	if err != nil {
		log.Error("decode err=%v", err)
		return
	}
	// log.Info("one=%v", one)
}

func TestJSON(t *testing.T) {
	now := time.Now()
	for i := 0; i < 10000; i++ {
		jsonload()
	}
	log.Info("json-------------------------cost=%v", time.Since(now))
}

var (
	count = 10
)

func mfunc() {
	count++
}

func TestHotfix(t *testing.T) {
	L := NewState(SetHotfix(true))
	defer L.Close()
	doErr := L.DoProFiles("./script/args.lua")
	if doErr != nil {
		log.Error("err=%v", doErr)
		return
	}
	L.RegMod("mymod", LStateMod{
		"func1": mfunc,
	})
	curr := time.Now()
	// for i := 0; i < 1000000; i++ {
	for i := 0; i < 10; i++ {
		_, err := L.Call("mymod.func1")
		if err != nil {
			log.Error("err=%v", err)
			return
		}
	}
	log.Info("cost=%v", time.Since(curr))
}
