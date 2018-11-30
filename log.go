package ilua

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	glua "github.com/yuin/gopher-lua"
)

// LogFuncType log function type
type LogFuncType func(string, ...interface{})

// LuaLog lua log
type LuaLog interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// SetLogger set logger
func SetLogger(log LuaLog) {
	logger = log
}

/////////////////////////////////////////////////////

// export log libs for lua

var logFuncs map[string]glua.LGFunction

// OpenLogLib export log lib
func OpenLogLib(L *glua.LState, modName string) {
	logFuncs = make(map[string]glua.LGFunction)
	logFuncs["Info"] = func(cb LogFuncType) glua.LGFunction {
		return func(L *glua.LState) int {
			return exportLog(L, cb)
		}
	}(loginfo)
	logFuncs["Error"] = func(cb LogFuncType) glua.LGFunction {
		return func(L *glua.LState) int {
			return exportLog(L, cb)
		}
	}(logerror)
	L.Push(L.NewFunction(exportLogLoader))
	L.Push(glua.LString(modName))
	L.Call(1, 0)
}

func exportLogLoader(L *glua.LState) int {
	modName, ok := L.Get(-1).(glua.LString)
	if !ok {
		logerror("log lib param error")
		return 0
	}
	mod := L.RegisterModule(string(modName), logFuncs).(*glua.LTable)
	L.Push(mod)
	return 1
}

func exportLog(L *glua.LState, logfunc LogFuncType) int {
	f, fok := L.Get(-2).(glua.LString)
	if !fok {
		logerror("format not string type")
		return 0
	}
	args, ok := L.Get(-1).(*glua.LTable)
	if !ok {
		logerror("args not table type")
		return 0
	}
	lastIdx := args.Len()
	if lastIdx <= 0 {
		return 0
	}
	lastIdx = lastIdx
	sizeVal, sizeOK := args.RawGet(glua.LNumber(lastIdx)).(glua.LNumber)
	if !sizeOK {
		logerror("args last one not number")
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
	return nil
}

/////////////////////////////////////////////////////

var logger LuaLog

func loginfo(format string, args ...interface{}) {
	logger.Info(format, args...)
}

func logerror(format string, args ...interface{}) {
	logger.Error(format, args...)
}

/////////////////////////////////////////////////////

type stdLog struct {
}

func (imp *stdLog) Info(format string, args ...interface{}) {
	imp.write("INFO", format, args...)
}

func (imp *stdLog) Error(format string, args ...interface{}) {
	imp.write("ERR ", format, args...)
}

func (imp *stdLog) write(level string, format string, args ...interface{}) {
	filename, line, funcname := "???", 0, "???"
	var ok bool
	var pc uintptr
	pc, filename, line, ok = runtime.Caller(3)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
		filename = filepath.Base(filename)
	}
	fstr := fmt.Sprintf("%s|%s|%d|%s:%d,%s: %s", genTime(), level, os.Getpid(), filename, line, funcname, fmt.Sprintf(format, args...))
	fmt.Println(fstr)
}

/////////////////////////////////////////////////////
// util

const (
	timeFormat = "20060102 15:04:05"
)

func genTime() string {
	return time.Now().Format(timeFormat)
}
