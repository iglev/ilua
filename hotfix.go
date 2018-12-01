package ilua

import (
	"time"

	"github.com/iglev/ilua/log"
	glua "github.com/yuin/gopher-lua"
)

func checkHotfix(L *LState) (err error) {
	curr := time.Now().Unix()
	log.Info("curr=%v last=%v hotfixtime=%v", curr, L.lastHotfixTime, L.opts.HotfixTime)
	if curr >= (L.lastHotfixTime + L.opts.HotfixTime) {
		err = L.L().CallByParam(glua.P{
			Fn:   L.L().GetGlobal(LuaFuncHotfix),
			NRet: 0,
		})
		if err != nil {
			log.Error("call hotfix err=%v", err)
		}
		L.lastHotfixTime = curr
	}
	return nil
}
