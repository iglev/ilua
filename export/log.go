package export

import (
	log "github.com/iglev/ilog"
	glua "github.com/yuin/gopher-lua"
)

const (
	// LogLibName log module name
	LogLibName = "mlog"
)

// LogFuncType log function type
type LogFuncType func(string, ...interface{})

// OpenLogLib export log lib
func OpenLogLib(L *glua.LState) {
	OpenLib(L, LogLibName, map[string]interface{}{
		"Info": func(cb LogFuncType) interface{} {
			return func(format string, args ...glua.LValue) {
				exportLog(cb, format, args...)
			}
		}(log.Info),
		"Error": func(cb LogFuncType) interface{} {
			return func(format string, args ...glua.LValue) {
				exportLog(cb, format, args...)
			}
		}(log.Error),
	})
}

func exportLog(logfunc LogFuncType, format string, args ...glua.LValue) {
	size := len(args)
	if size > 0 {
		fargs := make([]interface{}, 0, size)
		for i := 0; i < size; i++ {
			fargs = append(fargs, rawOut(args[i]))
		}
		logfunc(format, fargs...)
	} else {
		logfunc(format)
	}
}

func rawOut(value glua.LValue) interface{} {
	switch value.(type) {
	case *glua.LTable:
		tb := value.(*glua.LTable)
		return rawOutLTable(tb)
	case *glua.LUserData:
		ud := value.(*glua.LUserData)
		return rawOUtLUserData(ud)
	default:
		return value
	}
}

func rawOutLTable(tb *glua.LTable) interface{} {
	o := make([]interface{}, 0, 16)
	tb.ForEach(func(k, v glua.LValue) {
		o = append(o, rawOut(v))
	})
	return o
}

func rawOUtLUserData(ud *glua.LUserData) interface{} {
	return ud.Value
}
