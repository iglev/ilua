package ilua

import (
	"time"

	"github.com/iglev/ilua/log"
	glua "github.com/yuin/gopher-lua"
)

func checkHotfix(L *LState) (err error) {
	curr := time.Now().Unix()
	if (L.opts.HotfixTime >= 0) && curr >= (L.lastHotfixTime+L.opts.HotfixTime) {
		err = L.L().CallByParam(glua.P{
			Fn:      L.L().GetGlobal(LuaFuncHotfix),
			NRet:    0,
			Protect: true,
		})
		if err != nil {
			log.Error("call hotfix err=%v", err)
		}
		L.lastHotfixTime = curr
	}
	return nil
}
