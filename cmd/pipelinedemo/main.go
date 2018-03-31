package main

import (
	"bufio"
	"fmt"
	"os"

	"gointro/pipeline"
)

func main() {
	const filename = "small.in"
	const num = 64
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.RandomSource(num)
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, p)
	writer.Flush()

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(file, -1)

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
