package ilua

import (
	"errors"

	"github.com/iglev/ilua/log"
	glua "github.com/yuin/gopher-lua"
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

// doProFiles load main
func doProFiles(L *glua.LState, argsFile string) (err error) {
	err = L.DoFile(argsFile)
	if err != nil {
		log.Error("DoFile fail, argsFile=%v err=%v", argsFile, err)
		return
	}
	baseMainFile, bmOK := L.GetGlobal(LuaBaseMainFileName).(glua.LString)
	if !bmOK {
		log.Error("not found basefile=%v", LuaBaseMainFileName)
		err = ErrConfigScript
		return
	}
	err = loadLuaFiles(L, string(baseMainFile), LuaBaseMainModule)
	if err != nil {
		log.Error("loadLuaFiles fail, basemainfile=%v err=%v", string(baseMainFile), err)
		return
	}
	mainFile, mainOK := L.GetGlobal(LuaMainFileName).(glua.LString)
	if !mainOK {
		log.Error("not found mainfile=%v", LuaMainFileName)
		err = ErrConfigScript
		return
	}
	err = loadLuaFiles(L, string(mainFile), LuaMainModule)
	if err != nil {
		log.Error("loadLuaFiles fail, mainfile=%v err=%v", string(mainFile), err)
		return
	}
	return
}

func loadLuaFiles(L *glua.LState, mainfile, modName string) (err error) {
	err = L.DoFile(mainfile)
	if err != nil {
		log.Error("do mainfile fail, mainfile=%v err=%v", mainfile, err)
		return
	}
	mainMod, mainModOK := L.GetGlobal(modName).(*glua.LTable)
	if !mainModOK {
		log.Error("mainModule fail")
		return ErrMainModuleNotFound
	}
	luafiles, luafilesOK := L.GetField(mainMod, LuaValLuaFiles).(*glua.LTable)
	if !luafilesOK {
		log.Error("luafiles fail")
		return ErrLuaFilesNotFound
	}
	L.ForEach(luafiles, func(key glua.LValue, val glua.LValue) {
		submod, submodOK := val.(*glua.LTable)
		if !submodOK {
			log.Error("submod fail")
			err = ErrSubModNotTable
			return
		}
		dir, dirOK := L.GetField(submod, LuaValDir).(glua.LString)
		if !dirOK {
			log.Error("dir fail")
			err = ErrSubModDirError
			return
		}
		files, filesOK := L.GetField(submod, LuaValFiles).(*glua.LTable)
		if !filesOK {
			log.Error("files fail, dir=%v", string(dir))
			err = ErrSubModFilesError
			return
		}
		L.ForEach(files, func(fkey glua.LValue, fval glua.LValue) {
			filename, filenameOK := fval.(glua.LString)
			if !filenameOK {
				log.Error("filename fail, dir=%v val=%v", dir, val)
				err = ErrSubModFilenameError
				return
			}
			fullname := string(dir) + string(filename)
			err = L.DoFile(fullname)
			if err != nil {
				log.Error("DoFile fail, err=%v", err)
				return
			}
			log.Info("load success %v", string(filename))
		})
	})
	return
}
