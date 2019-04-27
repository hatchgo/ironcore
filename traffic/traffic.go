package traffic

import "sync"

// GoTraffic limit go concorrency level
type GoTraffic struct {
	mux sync.Mutex
	wg  sync.WaitGroup
	ch  chan int
	cnt int
}

// Add add one goroutine
func (traffic *GoTraffic) Add() int {
	traffic.ch <- 1
	traffic.mux.Lock()
	defer traffic.wg.Add(1)
	defer traffic.mux.Unlock()
	traffic.cnt++

	return traffic.cnt
}

// Done finish one goroutine
func (traffic *GoTraffic) Done() int {
	<-traffic.ch // channel cannot put inside locker
	traffic.mux.Lock()
	defer traffic.wg.Done()
	defer traffic.mux.Unlock()
	traffic.cnt--
	return traffic.cnt
}

// Wait for this traffic to Stop
func (traffic *GoTraffic) Wait() {
	traffic.wg.Wait()
}

// Control control max goroutine at a time if required
func Control(concurrency int) *GoTraffic {
	return &GoTraffic{
		ch: make(chan int, concurrency),
	}
}
