package main

import (
	"fmt"
	"sync"
)

func gen(done chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
		close(out)
	}()

	return out
}

func sq(done chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
		close(out)
	}()

	return out
}
func merge(done chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// Set pipeline
	nums := []int{2, 3, 4, 5, 6}
	done := make(chan struct{})
	defer close(done)

	in := gen(done, nums...)
	c1 := sq(done, in)
	c2 := sq(done, in)

	for n := range merge(done, c1, c2) {
		fmt.Println("demo3", n)
	}
}
