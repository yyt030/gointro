package main

import (
	"fmt"
	"os"

	"gointro/pipeline"
)

func main() {
	const filename = "largel.in"
	const num = 100000000
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.RandomSource(num)
	pipeline.WriterSink(file, p)

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(file)

	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count > 100 {
			break
		}
	}
}

func mergeDemo() {
	p := pipeline.Merge(
		pipeline.InMemSort(pipeline.ArraySource(3, 2, 16, 7, 4)),
		pipeline.InMemSort(pipeline.ArraySource(13, 2, 6, 17, 4)))
	for v := range p {
		fmt.Println(v)
	}
}
