package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []int{8, 1, 2, 9, 2, 1, 2, 3, 4}
	sort.Ints(a)
	for _, v := range a {
		fmt.Println(v)
	}
}
