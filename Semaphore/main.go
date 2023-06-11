package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	semaCh chan struct{}
}

func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.semaCh
}

func main() {
	var wg sync.WaitGroup
	semaphore := NewSemaphore(2)

	for idx := 0; idx < 10; idx++ {
		wg.Add(1)

		go func(taskID int) {
			semaphore.Acquire()

			defer wg.Done()
			defer semaphore.Release()

			msg := fmt.Sprintf(
				"%s Running worker %d",
				time.Now().Format("15:04:05"),
				taskID,
			)
			fmt.Println(msg)

			time.Sleep(1 * time.Second)
		}(idx)
	}

	wg.Wait()
}
