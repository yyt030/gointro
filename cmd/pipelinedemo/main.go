package main

import (
	"fmt"

	"gointro/pipeline"
)

func main() {
	p := pipeline.InMemSort(pipeline.ArraySource(3, 2, 6, 7, 4))
	for v := range p {
		fmt.Println(v)
	}

}
