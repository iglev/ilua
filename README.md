# ilua
<pre>
base on github.com/yuin/gopher-lua. 
1) 加载多个lua文件目录
2) 支持热更新
3) golang 和 lua 相互调用
4) 支持lua table(LTB) 反序列化成 golang struct对象
</pre>

# script 目录
<pre>
base: 基础libs
	- gfl: golang for lua (golang导出给lua接口)
	- lfg: lua for golang (lua提供给golang接口)
	- basemain.lua:  描述base目录下lua文件集合
pro: 逻辑模块
	- main.lua: 描述pro目录下lua文件集合 (业务按照模块方式管理)
args.lua: 入口参数,第一个加载lua文件(全局路径, 全局模块 和 path相关)
</pre>

# 导出golang struct

<pre>
type Person struct {
​	Name string
​	Age  int
​	BD   BirthDay
}

type BirthDay struct {
​	Y int16
​	M uint8
​	D uint8
}

func TestRegType(t *testing.T) {
​	L := NewState()
​	defer L.Close()
​	L.RegType("Person", Person{})
​	err := L.L().DoString(`
​		p = Person()
​		p.Name = "testname"
​		p.Age = 12
​		p.BD = { Y=2000, M=1, D=1 }
​		print(p.Name, p.Age, p.BD.Y, p.BD.M, p.BD.D)
​		print(type(p.BD)) -- userdata
​	`)
​	if err != nil {
​		log.Error("err=%v", err)
​		return
​	}
}
</pre>

# 导出golang函数 & 调用
<pre>
func createPerson(name string, age int) *Person {
	return &Person{Name: name, Age: age}
}

func TestCallFunc(t *testing.T) {
​	L := NewState()
​	defer L.Close()
​	doErr := L.DoProFiles("./script/args.lua")
​	if doErr != nil {
​		log.Error("err=%v", doErr)
​		return
​	}
​	L.RegMod("mymod", LStateMod{
​		"create": createPerson,
​		"incNum": 10,
​	})
​	ret, err := L.Call("mymod.create", "tname", 12)
​	if err != nil {
​		log.Error("err=%v", err)
​		return
​	}
​	res, resOK := ret.(*glua.LUserData)
​	if !resOK {
​		log.Error("not userdata, ret=%v", ret)
​		return
​	}
​	log.Info("res=%v", res.Value)
}
</pre>

# lua table 反序列成golang struct对象
<pre>
type ltb struct {
	Name     string
	Degree   bool
	LogLevel int
	Sub      subltb
	Va       int
}
type subltb struct {
​	N  string
​	XY subltb2
}

type subltb2 struct {
​	X, Y int64
}

func TestLTB(t *testing.T) {
​	L := NewState()
​	defer L.Close()
​	// res, _ := L.UnmarshalLTB("./script/conf.lua", ltb{})
​	res, _ := UnmarshalLTB("./script/conf.lua", ltb{})
​	log.Info("res=%+v", res)
}
</pre>
