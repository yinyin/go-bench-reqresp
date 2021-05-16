package reqresp

import (
	"sync"
)

type LockBasedReqResp struct {
	sync.Mutex

	Req  int32
	Resp int32
}

func (r *LockBasedReqResp) Wait(ch chan *LockBasedReqResp) {
	r.Lock()
	ch <- r
	r.Lock()
	defer r.Unlock()
}

func (r *LockBasedReqResp) Complete() {
	r.Unlock()
}
