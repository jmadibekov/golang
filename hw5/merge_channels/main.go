package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func channelGenerator(msg int) <-chan int {
	rand.Seed(time.Now().UnixNano())
	ch := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			ch <- msg
			time.Sleep(time.Duration(rand.Intn(1000) * int(time.Millisecond)))
		}

		close(ch)
	}()

	return ch
}

func merge(cs ...<-chan int) <-chan int {
	mergedCh := make(chan int)
	wg := new(sync.WaitGroup)

	for _, c := range cs {
		wg.Add(1)

		go func(localCh <-chan int) {
			defer wg.Done()

			for in := range localCh {
				mergedCh <- in
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(mergedCh)
	}()

	return mergedCh
}

func main() {
	mergedCh := merge(channelGenerator(0), channelGenerator(1), channelGenerator(3))

	for v := range mergedCh {
		fmt.Println(v)
	}
}
