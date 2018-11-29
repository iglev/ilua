package ilua

import (
	glua "github.com/yuin/gopher-lua"
)

// LoadMain load script
func LoadMain(L *glua.LState, configScript string) (err error) {
	err = L.DoFile(configScript)
	if err != nil {
		logerror("DoFile fail, configScript=%v err=%v", configScript, err)
		return
	}

	return
}
