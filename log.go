package ilua

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

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

var logger LuaLog

func loginfo(format string, args ...interface{}) {
	logger.Info(format, args...)
}

func logerror(format string, args ...interface{}) {
	logger.Error(format, args...)
}

/////////////////////////////////////////////////////

type stdLuaLogimp struct {
}

func (imp *stdLuaLogimp) Info(format string, args ...interface{}) {
	imp.write("INFO ", format, args...)
}

func (imp *stdLuaLogimp) Error(format string, args ...interface{}) {
	imp.write("ERROR", format, args...)
}

func (imp *stdLuaLogimp) write(level string, format string, args ...interface{}) {
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
