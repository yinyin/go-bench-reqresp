package reqresp_test

import (
	"sync"
	"testing"

	reqresp "github.com/yinyin/go-bench-reqresp"
)

func workerTestLockBasedReqRespPlus1(wg *sync.WaitGroup, ch chan *reqresp.LockBasedReqResp) {
	d := <-ch
	for d != nil {
		d.Resp = d.Req + 1
		d.Complete()
		// prepare next iteration
		d = <-ch
	}
	wg.Done()
}

func TestLockBasedReqResp_WaitComplete_1pass(t *testing.T) {
	target := reqresp.LockBasedReqResp{
		Req: 7,
	}
	ch := make(chan *reqresp.LockBasedReqResp)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go (func() {
		d := <-ch
		d.Resp = d.Req + 1
		d.Complete()
		wg.Done()
	})()
	target.Wait(ch)
	if target.Resp != 8 {
		t.Errorf("unexpect response: %d (req=%d)", target.Resp, target.Req)
	}
	close(ch)
	wg.Wait()
}

func TestLockBasedReqResp_WaitComplete_90pass(t *testing.T) {
	target := reqresp.LockBasedReqResp{
		Req: 7,
	}
	ch := make(chan *reqresp.LockBasedReqResp)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go workerTestLockBasedReqRespPlus1(&wg, ch)
	for attempt := int32(0); attempt < 90; attempt++ {
		target.Req = attempt + 7
		target.Resp = 0
		target.Wait(ch)
		if target.Resp != attempt+8 {
			t.Errorf("unexpect response(attempt=%d): %d (req=%d)", attempt, target.Resp, target.Req)
		}
	}
	close(ch)
	wg.Wait()
}

func implBenchmarkParallelLockBasedReqResp_WaitComplete(b *testing.B, workerCount int, workBufferSize int) {
	wg := sync.WaitGroup{}
	ch := make(chan *reqresp.LockBasedReqResp, workBufferSize)
	for workerIdx := 0; workerIdx < workerCount; workerIdx++ {
		wg.Add(1)
		go workerTestLockBasedReqRespPlus1(&wg, ch)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		target := reqresp.LockBasedReqResp{}
		reqVal := int32(3)
		for pb.Next() {
			target.Req = reqVal
			target.Resp = 0
			target.Wait(ch)
			_ = target.Resp
			reqVal++
		}
	})
	b.StopTimer()
	close(ch)
	wg.Wait()
}

func BenchmarkParallelLockBasedReqResp_WaitComplete_1_1(b *testing.B) {
	implBenchmarkParallelLockBasedReqResp_WaitComplete(b, 1, 1)
}

func BenchmarkParallelLockBasedReqResp_WaitComplete_2_1(b *testing.B) {
	implBenchmarkParallelLockBasedReqResp_WaitComplete(b, 2, 1)
}

func BenchmarkParallelLockBasedReqResp_WaitComplete_2_2(b *testing.B) {
	implBenchmarkParallelLockBasedReqResp_WaitComplete(b, 2, 2)
}

func BenchmarkParallelLockBasedReqResp_WaitComplete_2_4(b *testing.B) {
	implBenchmarkParallelLockBasedReqResp_WaitComplete(b, 2, 4)
}

func BenchmarkParallelLockBasedReqResp_WaitComplete_4_2(b *testing.B) {
	implBenchmarkParallelLockBasedReqResp_WaitComplete(b, 4, 2)
}
