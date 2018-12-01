package export

import (
	"github.com/iglev/ilua/log"
	glua "github.com/yuin/gopher-lua"
)

const (
	// LogLibName log module name
	LogLibName = "MLog"
)

// LogFuncType log function type
type LogFuncType func(string, ...interface{})

// OpenLogLib export log lib
func OpenLogLib(L *glua.LState) {
	OpenLib(L, LogLibName, map[string]glua.LGFunction{
		"Info": func(cb LogFuncType) glua.LGFunction {
			return func(L *glua.LState) int {
				return exportLog(L, cb)
			}
		}(log.Info),
		"Error": func(cb LogFuncType) glua.LGFunction {
			return func(L *glua.LState) int {
				return exportLog(L, cb)
			}
		}(log.Error),
	})
}

func exportLog(L *glua.LState, logfunc LogFuncType) int {
	f, fok := L.Get(-2).(glua.LString)
	if !fok {
		log.Error("format not string type")
		return 0
	}
	args, ok := L.Get(-1).(*glua.LTable)
	if !ok {
		log.Error("args not table type")
		return 0
	}
	lastIdx := args.Len()
	if lastIdx <= 0 {
		return 0
	}
	sizeVal, sizeOK := args.RawGet(glua.LNumber(lastIdx)).(glua.LNumber)
	if !sizeOK {
		log.Error("args last one not number")
		return 0
	}
	size := int(sizeVal)
	if size > 0 {
		fargs := make([]interface{}, 0, size)
		for i := 1; i <= size; i++ {
			fargs = append(fargs, rawOut(args.RawGet(glua.LNumber(i))))
		}
		logfunc(string(f), fargs...)
	} else {
		logfunc(string(f))
	}
	return 0
}

func rawOut(value glua.LValue) interface{} {
	switch value.(type) {
	case *glua.LTable:
		tb := value.(*glua.LTable)
		o := make([]interface{}, 0, 16)
		tb.ForEach(func(k, v glua.LValue) {
			o = append(o, rawOut(v))
		})
		return o
	default:
		return value
	}
}
