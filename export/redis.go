package export

import (
	"github.com/gomodule/redigo/redis"
	glua "github.com/yuin/gopher-lua"
)

const (
	// RedisLibName redis module name
	RedisLibName = "redis"
)

// OpenRedisLib export redis lib
func OpenRedisLib(L *glua.LState) {
	mod := OpenLib(L, RedisLibName, map[string]interface{}{
		// reply
		"Int":        rInt,
		"Int64":      rInt64,
		"Uint64":     rUint64,
		"Float64":    rFloat64,
		"String":     rString,
		"Bytes":      rBytes,
		"Bool":       rBool,
		"Values":     rValues,
		"Float64s":   rFloat64s,
		"Strings":    rStrings,
		"ByteSlices": rByteSlices,
		"Int64s":     rInt64s,
		"Ints":       rInts,
		"StringMap":  rStringMap,
		"IntMap":     rIntMap,
		"Int64Map":   rInt64Map,
		"Positions":  rPositions,
		// var
		"ErrNil": redis.ErrNil,
	})
	// type
	NewModType(L, mod, "Args", redis.Args{})
}

func rInt(conn redis.Conn, commandName string, args ...interface{}) (int, error) {
	return redis.Int(conn.Do(commandName, args...))
}

func rInt64(conn redis.Conn, commandName string, args ...interface{}) (int64, error) {
	return redis.Int64(conn.Do(commandName, args...))
}

func rUint64(conn redis.Conn, commandName string, args ...interface{}) (uint64, error) {
	return redis.Uint64(conn.Do(commandName, args...))
}

func rFloat64(conn redis.Conn, commandName string, args ...interface{}) (float64, error) {
	return redis.Float64(conn.Do(commandName, args...))
}

func rString(conn redis.Conn, commandName string, args ...interface{}) (string, error) {
	return redis.String(conn.Do(commandName, args...))
}

func rBytes(conn redis.Conn, commandName string, args ...interface{}) ([]byte, error) {
	return redis.Bytes(conn.Do(commandName, args...))
}

func rBool(conn redis.Conn, commandName string, args ...interface{}) (bool, error) {
	return redis.Bool(conn.Do(commandName, args...))
}

func rValues(conn redis.Conn, commandName string, args ...interface{}) ([]interface{}, error) {
	return redis.Values(conn.Do(commandName, args...))
}

func rFloat64s(conn redis.Conn, commandName string, args ...interface{}) ([]float64, error) {
	return redis.Float64s(conn.Do(commandName, args...))
}

func rStrings(conn redis.Conn, commandName string, args ...interface{}) ([]string, error) {
	return redis.Strings(conn.Do(commandName, args...))
}

func rByteSlices(conn redis.Conn, commandName string, args ...interface{}) ([][]byte, error) {
	return redis.ByteSlices(conn.Do(commandName, args...))
}

func rInt64s(conn redis.Conn, commandName string, args ...interface{}) ([]int64, error) {
	return redis.Int64s(conn.Do(commandName, args...))
}

func rInts(conn redis.Conn, commandName string, args ...interface{}) ([]int, error) {
	return redis.Ints(conn.Do(commandName, args...))
}

func rStringMap(conn redis.Conn, commandName string, args ...interface{}) (map[string]string, error) {
	return redis.StringMap(conn.Do(commandName, args...))
}

func rIntMap(conn redis.Conn, commandName string, args ...interface{}) (map[string]int, error) {
	return redis.IntMap(conn.Do(commandName, args...))
}

func rInt64Map(conn redis.Conn, commandName string, args ...interface{}) (map[string]int64, error) {
	return redis.Int64Map(conn.Do(commandName, args...))
}

func rPositions(conn redis.Conn, commandName string, args ...interface{}) ([]*[2]float64, error) {
	return redis.Positions(conn.Do(commandName, args...))
}
