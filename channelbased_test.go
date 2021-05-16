package reqresp_test

import (
	"sync"
	"testing"

	reqresp "github.com/yinyin/go-bench-reqresp"
)

func workerTestChannelBasedReqRespPlus1(wg *sync.WaitGroup, ch chan *reqresp.ChannelBasedReqResp) {
	d := <-ch
	for d != nil {
		d.Complete(d.Req + 1)
		// prepare next iteration
		d = <-ch
	}
	wg.Done()
}

func TestChannelBasedReqResp_WaitComplete_1pass(t *testing.T) {
	target := reqresp.NewChannelBasedReqResp(7)
	ch := make(chan *reqresp.ChannelBasedReqResp)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go (func() {
		d := <-ch
		d.Complete(d.Req + 1)
		wg.Done()
	})()
	target.Wait(ch)
	if respVal := <-target.RespCh; respVal != 8 {
		t.Errorf("unexpect response: %d (req=%d)", respVal, target.Req)
	}
	close(ch)
	wg.Wait()
}

func TestChannelBasedReqResp_WaitComplete_90pass(t *testing.T) {
	target := reqresp.NewChannelBasedReqResp(7)
	ch := make(chan *reqresp.ChannelBasedReqResp)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go workerTestChannelBasedReqRespPlus1(&wg, ch)
	for attempt := int32(0); attempt < 90; attempt++ {
		target.Req = attempt + 7
		target.Wait(ch)
		if respVal := <-target.RespCh; respVal != attempt+8 {
			t.Errorf("unexpect response(attempt=%d): %d (req=%d)", attempt, respVal, target.Req)
		}
	}
	close(ch)
	wg.Wait()
}

func implBenchmarkParallelChannelBasedReqResp_WaitComplete(b *testing.B, workerCount int, workBufferSize int) {
	wg := sync.WaitGroup{}
	ch := make(chan *reqresp.ChannelBasedReqResp, workBufferSize)
	for workerIdx := 0; workerIdx < workerCount; workerIdx++ {
		wg.Add(1)
		go workerTestChannelBasedReqRespPlus1(&wg, ch)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		target := reqresp.NewChannelBasedReqResp(0)
		reqVal := int32(3)
		for pb.Next() {
			target.Req = reqVal
			target.Wait(ch)
			respVal := <-target.RespCh
			_ = respVal
			reqVal++
		}
	})
	b.StopTimer()
	close(ch)
	wg.Wait()
}

func BenchmarkParallelChannelBasedReqResp_WaitComplete_1_1(b *testing.B) {
	implBenchmarkParallelChannelBasedReqResp_WaitComplete(b, 1, 1)
}

func BenchmarkParallelChannelBasedReqResp_WaitComplete_2_1(b *testing.B) {
	implBenchmarkParallelChannelBasedReqResp_WaitComplete(b, 2, 1)
}

func BenchmarkParallelChannelBasedReqResp_WaitComplete_2_2(b *testing.B) {
	implBenchmarkParallelChannelBasedReqResp_WaitComplete(b, 2, 2)
}

func BenchmarkParallelChannelBasedReqResp_WaitComplete_2_4(b *testing.B) {
	implBenchmarkParallelChannelBasedReqResp_WaitComplete(b, 2, 4)
}

func BenchmarkParallelChannelBasedReqResp_WaitComplete_4_2(b *testing.B) {
	implBenchmarkParallelChannelBasedReqResp_WaitComplete(b, 4, 2)
}
