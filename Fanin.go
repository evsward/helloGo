package main

import (
	"sync"
	"fmt"
	"hello/pipeline"
)

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	// 定义一个独立通道out，将接收所有值
	out := make(chan int)

	// 将通道中所有值转到out
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
		wg.Done()
	}

	wg.Add(len(cs))
	// 将merge参数中所有通道的值都合到唯一通道out上去
	for _, c := range cs {
		go output(c)
	}

	// 启动一个额外的goroutine（不会按照代码顺序执行，而是一进到merge就会启动）来等待直到所有通道Done以后关闭那个唯一通道out。
	go func() {
		wg.Wait() // 直到wg全都Done了才会继续执行。
		close(out)
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)
	in := pipeline.Gen(2, 3)
	c1 := pipeline.Sq(in)
	c2 := pipeline.Sq(in)
	out := merge(done, c1, c2)
	fmt.Println(<-out)
}

// Output:
//	4
