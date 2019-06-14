package utils

import (
	"sync"
)

type IntPipe func(<-chan int) <-chan int

func Chain(ps ...IntPipe) IntPipe {
	return func(in <-chan int) <-chan int {
		c := in
		for _, p := range ps {
			c = p(c)
		}
		return c
	}
}

func FanIn(ins ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(ins))

	for _, in := range ins {
		go func(in <-chan int) {
			defer wg.Done()
			for n := range in {
				out <- n
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Distribution(p IntPipe, n int) IntPipe {
	return func(in <-chan int) <-chan int {
		cs := make([]<-chan int, n)
		
		for i := 0; i < n ; i++ {
			cs[i] = p(in)
		}

		return FanIn(cs...)
	}
}
