package ilua

import (
	"errors"

	glua "github.com/yuin/gopher-lua"
)

const (
	luaMainFileName = "G_MAIN_FILE"
	luaMainModule   = "G_MAIN"
	luaValLuaFiles  = "luafiles"
	luaValDir       = "dir"
	luaValFiles     = "files"
	luaValDecode    = "decode"
)

var (
	// ErrConfigScript config script error
	ErrConfigScript = errors.New("configScript error")
	// ErrMainModuleNotFound main module not found error
	ErrMainModuleNotFound = errors.New("main module not found error")
	// ErrLuaFilesNotFound luafiles not found error
	ErrLuaFilesNotFound = errors.New("luafiles not found error")
	// ErrSubModNotTable submodule not table error
	ErrSubModNotTable = errors.New("submodule not table error")
	// ErrSubModDirError submodule dir error
	ErrSubModDirError = errors.New("submodule dir error")
	// ErrSubModFilesError submodule files error
	ErrSubModFilesError = errors.New("submodule files error")
	// ErrSubModFilenameError submodule filename error
	ErrSubModFilenameError = errors.New("submodule filename error")
)

// LoadLibs load main
func LoadLibs(L *glua.LState, configScript string) (err error) {
	err = L.DoFile(configScript)
	if err != nil {
		logerror("DoFile fail, configScript=%v err=%v", configScript, err)
		return
	}
	mainPath, mainOK := L.GetGlobal(luaMainFileName).(glua.LString)
	if !mainOK {
		logerror("not found mainfile=%v", luaMainFileName)
		err = ErrConfigScript
		return
	}
	mainfile := string(mainPath)
	loginfo("mainfile=%v", mainfile)
	err = loadLuaFiles(L, mainfile)
	if err != nil {
		logerror("loadLuaFiles fail, mainfile=%v err=%v", mainfile, err)
		return
	}
	return
}

func loadLuaFiles(L *glua.LState, mainfile string) (err error) {
	err = L.DoFile(mainfile)
	if err != nil {
		logerror("do mainfile fail, mainfile=%v err=%v", mainfile, err)
		return
	}
	mainMod, mainModOK := L.GetGlobal(luaMainModule).(*glua.LTable)
	if !mainModOK {
		logerror("mainModule fail")
		return ErrMainModuleNotFound
	}
	luafiles, luafilesOK := L.GetField(mainMod, luaValLuaFiles).(*glua.LTable)
	if !luafilesOK {
		logerror("luafiles fail")
		return ErrLuaFilesNotFound
	}
	L.ForEach(luafiles, func(key glua.LValue, val glua.LValue) {
		submod, submodOK := val.(*glua.LTable)
		if !submodOK {
			logerror("submod fail")
			err = ErrSubModNotTable
			return
		}
		dir, dirOK := L.GetField(submod, luaValDir).(glua.LString)
		if !dirOK {
			logerror("dir fail")
			err = ErrSubModDirError
			return
		}
		files, filesOK := L.GetField(submod, luaValFiles).(*glua.LTable)
		if !filesOK {
			logerror("files fail, dir=%v", string(dir))
			err = ErrSubModFilesError
			return
		}
		L.ForEach(files, func(fkey glua.LValue, fval glua.LValue) {
			filename, filenameOK := fval.(glua.LString)
			if !filenameOK {
				logerror("filename fail, dir=%v val=%v", dir, val)
				err = ErrSubModFilenameError
				return
			}
			fullname := string(dir) + string(filename)
			err = L.DoFile(fullname)
			if err != nil {
				logerror("DoFile fail, err=%v", err)
				return
			}
		})
	})
	return
}
