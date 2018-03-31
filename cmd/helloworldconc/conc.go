package main

import (
	"fmt"
)

func printHelloWorld(i int, ch chan string) {
	for {
		ch <- fmt.Sprintf("hello world from goroutine %d\n", i)
	}

}

func main() {
	ch := make(chan string)
	for i := 0; i < 1000; i++ {
		go printHelloWorld(i, ch)
	}

	for {
		msg := <-ch
		fmt.Println(msg)
	}
}
