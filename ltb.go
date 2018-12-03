package ilua

import (
	"errors"
	"fmt"
	"reflect"

	glua "github.com/yuin/gopher-lua"
)

var (
	// ErrInvalidUnmarshalVal - unmarshal val need ptr type error
	ErrInvalidUnmarshalVal = errors.New("unmarshal val need ptr type")
)

type ltbDecodeState struct {
	L   *glua.LState
	ltb *glua.LTable
}

func unmarshalLTB(L *glua.LState, script string, ival interface{}) error {
	err := L.DoFile(script)
	if err != nil {
		return err
	}
	tb, ok := L.Get(-1).(*glua.LTable)
	if !ok {
		return fmt.Errorf("script=%v invalid, `return { key = val, ... }`", script)
	}
	ds := &ltbDecodeState{
		L:   L,
		ltb: tb,
	}
	return ds.unmarshal(ival)
}

func (s *ltbDecodeState) unmarshal(ival interface{}) error {
	rv := reflect.ValueOf(ival)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidUnmarshalVal
	}
	// TODO ...
	return nil
}
