package main

import (
	"fmt"

	"gointro/pipeline"
)

func main() {
	p := pipeline.Merge(
		pipeline.InMemSort(pipeline.ArraySource(3, 2, 16, 7, 4)),
		pipeline.InMemSort(pipeline.ArraySource(13, 2, 6, 17, 4)))

	for v := range p {
		fmt.Println(v)
	}
}
