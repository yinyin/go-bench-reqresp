package reqresp

type ChannelBasedReqResp struct {
	RespCh chan int32

	Req int32
}

func NewChannelBasedReqResp(req int32) (r *ChannelBasedReqResp) {
	return &ChannelBasedReqResp{
		RespCh: make(chan int32, 1),
		Req:    req,
	}
}

func (r *ChannelBasedReqResp) Wait(ch chan *ChannelBasedReqResp) {
	ch <- r
}

func (r *ChannelBasedReqResp) Complete(resp int32) {
	r.RespCh <- resp
}

func (r *ChannelBasedReqResp) Release() {
	close(r.RespCh)
}
