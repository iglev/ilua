package ilua

import (
	"context"
	"sync"
	"time"
	"sync/atomic"

	log "github.com/iglev/ilog"
)

type hotfixMgr interface {
	reg(file string)
	check(L *LState) *hotfixList
}

func newHotfixMgr(ctx context.Context, needCoro bool, hotfixTime int64) hotfixMgr {
	var ht hotfixMgr
	if needCoro {
		htco := &hotfixMgrCoro{
			ctx:        ctx,
			ch:         make(chan *hotfixList, 20480),
			hotfixTime: hotfixTime,
		}
		go htco.loop()
		ht = htco
	} else {
		ht = &hotfixMgrLocal{
			mp:         make(map[string]int64),
			lasttime:   time.Now().Unix(),
			hotfixTime: hotfixTime,
		}
	}
	return ht
}

type hotfixList struct {
	files []string
}

func hotfixDoFile(L *LState, ht hotfixMgr, up *hotfixList) {
	size := len(up.files)
	for i := 0; i < size; i++ {
		err := L.L().DoFile(up.files[i])
		if err != nil {
			log.Error("DoFile fail, file=%v err=%v", up.files[i], err)
		} else {
			log.Info("reload file=%v success", up.files[i])
		}
	}
}

//////////////////////////////////////////////////////////////////
// hotfixMgrLocal

type hotfixMgrLocal struct {
	mp         map[string]int64
	lasttime   int64
	hotfixTime int64
}

func (ht *hotfixMgrLocal) reg(file string) {
	mt, err := getFileModtime(file)
	if err != nil {
		return
	}
	ht.mp[file] = mt
}

func (ht *hotfixMgrLocal) check(L *LState) *hotfixList {
	curr := time.Now().Unix()
	lasttime := atomic.LoadInt64(&ht.lasttime)
	if curr < (lasttime + ht.hotfixTime) {
		return nil
	}
	atomic.StoreInt64(&ht.lasttime, curr)
	up := ht.getHotfixList()
	if up != nil {
		hotfixDoFile(L, ht, up)
	}
	return up
}

func (ht *hotfixMgrLocal) getHotfixList() *hotfixList {
	size := len(ht.mp)
	if size <= 0 {
		return nil
	}
	up := &hotfixList{
		files: make([]string, 0, size),
	}
	tmp := make(map[string]int64)
	for k, v := range ht.mp {
		mt, err := getFileModtime(k)
		if err != nil {
			log.Error("getFileModtime fail, file=%v err=%v", k, err)
			continue
		}
		if v != mt {
			up.files = append(up.files, k)
			tmp[k] = mt
		}
	}
	for k, v := range tmp {
		ht.mp[k] = v
	}
	return up
}

//////////////////////////////////////////////////////////////////
// hotfixMgrCoro

type hotfixMgrCoro struct {
	ctx        context.Context
	mp         sync.Map
	ch         chan *hotfixList
	hotfixTime int64
}

func (ht *hotfixMgrCoro) reg(file string) {
	mt, err := getFileModtime(file)
	if err != nil {
		return
	}
	ht.mp.Store(file, mt)
}

func (ht *hotfixMgrCoro) check(L *LState) *hotfixList {
	up := ht.getHotfixList()
	if up != nil {
		hotfixDoFile(L, ht, up)
	}
	return up
}

func (ht *hotfixMgrCoro) getHotfixList() *hotfixList {
	select {
	case up := <-ht.ch:
		return up
	default:
		return nil
	}
}

func (ht *hotfixMgrCoro) loop() {
	timer := time.NewTimer(time.Duration(ht.hotfixTime) * time.Second)
	defer timer.Stop()
Loop:
	for {
		select {
		case <-timer.C:
			ht.loopCheck()
			timer.Reset(time.Duration(ht.hotfixTime) * time.Second)
		case <-ht.ctx.Done():
			break Loop
		}
	}
}

func (ht *hotfixMgrCoro) loopCheck() {
	var up *hotfixList
	mp := make(map[string]int64)
	ht.mp.Range(func(k, v interface{}) bool {
		file, ok := k.(string)
		if !ok {
			return false
		}
		oldmt, mtOK := v.(int64)
		if !mtOK {
			return false
		}
		newmt, err := getFileModtime(file)
		if err != nil {
			log.Error("getFileModtime fail, file=%v err=%v", file, err)
			return false
		}
		if newmt != oldmt {
			mp[file] = newmt
		}
		return true
	})
	size := len(mp)
	if size > 0 {
		up = &hotfixList{
			files: make([]string, 0, size),
		}
		for k, v := range mp {
			ht.mp.Store(k, v)
			up.files = append(up.files, k)
		}
		select {
		case ht.ch <- up:
			return
		default:
			log.Error("hotfix send to channel fail, up=%v", up)
		}
	}
}
