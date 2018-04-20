package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()

	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()

	return out
}
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
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
	c := gen(nums...)
	out := sq(c)

	for n := range out {
		fmt.Println("demo1", n)
	}

	for n := range sq(gen(nums...)) {
		fmt.Println("demo2", n)
	}

	c1 := sq(gen(nums...))
	c2 := sq(gen(nums...))

	for n := range merge(c1, c2) {
		fmt.Println("demo3", n)
	}
}
