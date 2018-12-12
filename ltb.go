package ilua

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/iglev/ilua/export"
	"github.com/iglev/ilua/log"
	glua "github.com/yuin/gopher-lua"
)

var (
	// ErrValueInvalid - value invalid error
	ErrValueInvalid = errors.New("value must `reflect.Struct` type error")
	// ErrUnmarshal - unmarshal error
	ErrUnmarshal     = errors.New("unmarshal error")
	decodeLuaFuncPre = "unmarshal"
	decodeLuaTypeTag = "#regtype#"
	decodeLuaCode    = `
		function unmarshal#regtype#(luatab)
			local tb = #regtype#()
			for k, v in pairs(luatab) do
				tb[k] = v
			end
			return tb
		end
	`
)

func unmarshalLTB(L *glua.LState, script string, ival interface{}) (interface{}, error) {
	err := L.DoFile(script)
	if err != nil {
		return nil, err
	}
	tb, ok := L.Get(-1).(*glua.LTable)
	if !ok {
		return nil, fmt.Errorf("script=%v invalid, `return { key = val, ... }`", script)
	}
	L.Pop(1)
	// reg type
	rv := reflect.ValueOf(ival)
	if rv.Kind() != reflect.Struct {
		return nil, ErrValueInvalid
	}
	k := rv.Type().String()
	ks := strings.Split(k, ".")
	regType := ks[len(ks)-1]
	lfuncStr := decodeLuaFuncPre + regType
	lfunc := L.GetGlobal(lfuncStr)
	if lfunc == glua.LNil {
		export.NewType(L, regType, ival)
		luastr := strings.Replace(decodeLuaCode, decodeLuaTypeTag, regType, -1)
		err = L.DoString(luastr)
		if err != nil {
			log.Error("dostring fail, str=%v err=%v", luastr, err)
			return nil, err
		}
		lfunc = L.GetGlobal(lfuncStr)
	}
	// exec unmarshal
	err = L.CallByParam(glua.P{
		Fn:      lfunc,
		NRet:    1,
		Protect: true,
	}, tb)
	if err != nil {
		log.Error("unmarshal fail, str=%v err=%v", lfuncStr, err)
		return nil, err
	}
	ud, udOK := L.Get(-1).(*glua.LUserData)
	if !udOK {
		log.Error("unmarshal get res fail, str=%v", lfuncStr)
		return nil, ErrUnmarshal
	}
	L.Pop(1)
	return ud.Value, nil
}
