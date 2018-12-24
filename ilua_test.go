package ilua

import (
	"sync/atomic"
	"testing"
	"time"

	log "github.com/iglev/ilog"
	glua "github.com/yuin/gopher-lua"
)

//////////////////////////////////////////////////////////////////////////////
// example: export go struct

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

//////////////////////////////////////////////////////////////////////////////
// example: reg module (function, value) & exec call

func TestCallFunc(t *testing.T) {
	L := NewState()
	defer L.Close()
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

//////////////////////////////////////////////////////////////////////////////
// example: lua table to go struct

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

func TestLTB(t *testing.T) {
	L := NewState()
	defer L.Close()
	res, _ := L.UnmarshalLTB("./script/conf.lua", ltb{})
	// res, _ := UnmarshalLTB("./script/conf.lua", ltb{})
	log.Info("res=%v", res)
}

//////////////////////////////////////////////////////////////////////////////
// example: start hotfix

func TestHotfix(t *testing.T) {
	// L := NewState(SetHofix(false)) // when call check file
	L := NewState(SetHotfix(true)) // go coroutine check file
	defer L.Close()
}

//////////////////////////////////////////////////////////////////////////////
// example: open redis lib

var (
	redisfunc = `function redisFunc(client)
					cnt, err = redis.Int(client, "HGET", "hash:mh", "num")
					if err ~= nil and err ~= redis.ErrNil then
						LogError("err=%v", err)
						return 0, err
					end
					return cnt, nil
				end`
)

func TestRedis(t *testing.T) {
	/*
		L := NewState()
		defer L.Close()
		export.OpenRedisLib(L.L())
		err := L.DoString(redisfunc)
		if err != nil {
			log.Error("err=%v", err)
			return
		}
		client, cerr := redis.Dial("tcp", "127.0.0.1:6379")
		if cerr != nil {
			log.Error("cerr=%v", cerr)
		}
		defer client.Close()
		res, resErr := L.Call("redisFunc", client)
		if resErr != nil {
			log.Error("err=%v", resErr)
			return
		}
		log.Info("res=%v", res)
	*/
}

func TestPool(t *testing.T) {
	var count int32
	pool := NewLStatePool(10, func() *LState {
		L := NewState()
		L.RegMod("mymod", LStateMod{
			"func": func() {
				atomic.AddInt32(&count, 1)
			},
		})
		return L
	})
	pool.Init()
	log.Info("pool=%v", pool)
	for i := 0; i < 100000; i++ {
		go func() {
			L := pool.Get()
			if L == nil {
				log.Error("pool had been closed")
				return
			}
			defer L.Close()
			L.Call("mymod.func")
		}()
	}
	time.Sleep(time.Second)
	pool.Close()
	log.Info("success, cnt=%v", count)
}

func TestLog(t *testing.T) {
	log.Info("1111111")
}
