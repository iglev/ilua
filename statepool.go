package ilua

import (
	"sync/atomic"
	"time"

	"github.com/iglev/ilua/log"
)

// LStatePool lua state pool
type LStatePool struct {
	ch        chan *LState
	close     int32
	size      uint32
	closesize uint32
	create    NewLStateFuncType
}

// NewLStateFuncType new lua state func type
type NewLStateFuncType func() *LState

// NewLStatePool new lua state pool
func NewLStatePool(poolSize int, createFunc NewLStateFuncType) *LStatePool {
	return &LStatePool{
		ch:     make(chan *LState, poolSize),
		close:  1,
		size:   (uint32)(poolSize),
		create: createFunc,
	}
}

// Init init pool
func (p *LStatePool) Init() error {
	size := (int)(p.size)
	for i := 0; i < size; i++ {
		L := p.create()
		L.linkPool(p)
		p.ch <- L
	}
	return nil
}

// Close close lua state pool
func (p *LStatePool) Close() {
	atomic.StoreInt32(&p.close, 0)
	// go func() {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
Loop:
	for {
		select {
		case L := <-p.ch:
			p.realClose(L)
			if atomic.LoadUint32(&p.closesize) == p.size {
				break Loop
			}
		case <-timer.C:
			if atomic.LoadUint32(&p.closesize) == p.size {
				break Loop
			}
			timer.Reset(time.Second)
		}
	}
	// }()
}

// Get get lua state
func (p *LStatePool) Get() *LState {
	if atomic.LoadInt32(&p.close) == 0 {
		return nil
	}
	return <-p.ch
}

func (p *LStatePool) pushback(L *LState) {
	if atomic.LoadInt32(&p.close) == 0 {
		p.realClose(L)
		return
	}
	select {
	case p.ch <- L:
	default:
		L.detachPool()
		L.Close()
		log.Error("pool error, channel push back fail!!!")
	}
}

func (p *LStatePool) realClose(L *LState) {
	L.detachPool()
	L.Close()
	atomic.AddUint32(&p.closesize, 1)
}
