package main

import (
	"sync/atomic"
	"time"
)

type Task func() string

type FutureResult struct {
	Done       atomic.Bool
	ResultChan chan string
}

func Async(t Task) *FutureResult {
	f := &FutureResult{
		ResultChan: make(chan string, 1),
	}
	go func() {
		res := t()
		f.ResultChan <- res
		f.Done.Store(true)
	}()
	return f
}

func AsyncWithTimeout(t Task, timeout time.Duration) *FutureResult {
	f := &FutureResult{
		ResultChan: make(chan string, 1),
	}

	go func() {
		done := make(chan string, 1)
		go func() {
			res := t()
			done <- res
		}()

		select {
		case res := <-done:
			f.ResultChan <- res
			f.Done.Store(true)
		case <-time.After(timeout):
			f.ResultChan <- "timeout"
		}
	}()

	return f
}

func (fResult *FutureResult) Await() string {
	return <-fResult.ResultChan
}

func CombineFutureResults(fResults ...*FutureResult) *FutureResult {
	total := len(fResults)
	f := &FutureResult{
		ResultChan: make(chan string, total),
	}

	go func() {
		for _, fr := range fResults {
			res := <-fr.ResultChan
			f.ResultChan <- res
		}
	}()

	return f
}

func exampleTask() string {
	time.Sleep(1 * time.Second)
	return "result"
}

func main() {
	f1 := Async(exampleTask)
	f2 := Async(exampleTask)

	results := CombineFutureResults(f1, f2)

	for i := 0; i < 2; i++ {
		res := <-results.ResultChan
		println(res)
	}
}
