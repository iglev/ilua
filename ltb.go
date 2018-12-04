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
	decodeLuaFunc    = "unmarshalLTB"
	decodeLuaTypeTag = "#regtype#"
	decodeLuaCode    = `
		function unmarshalLTB(luatab)
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
	export.NewType(L, regType, ival)
	// load unmarshal func
	luastr := strings.Replace(decodeLuaCode, decodeLuaTypeTag, regType, 1)
	err = L.DoString(luastr)
	if err != nil {
		log.Error("dostring fail, str=%v err=%v", luastr, err)
		return nil, err
	}
	// exec unmarshal
	err = L.CallByParam(glua.P{
		Fn:      L.GetGlobal(decodeLuaFunc),
		NRet:    1,
		Protect: true,
	}, tb)
	if err != nil {
		log.Error("unmarshal fail, str=%v err=%v", luastr, err)
		return nil, err
	}
	ud, udOK := L.Get(-1).(*glua.LUserData)
	if !udOK {
		log.Error("unmarshal get res fail, str=%v", luastr)
		return nil, ErrUnmarshal
	}
	L.Pop(1)
	return ud.Value, nil
}
