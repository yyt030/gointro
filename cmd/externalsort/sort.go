package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"gointro/pipeline"
)

func main() {
	const filename = "large.in"

	p := createNetworkPipeline(filename, 800000000, 4)
	writeToFile(p, "large.out")
	printFile("large.out")
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count > 100 {
			break
		}
	}
}

func writeToFile(c <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, c)

}

func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	pipeline.Init()
	chunkSize := fileSize / chunkCount
	sortResults := []<-chan int{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)

		sortResults = append(sortResults, pipeline.InMemSort(source))

	}
	return pipeline.MergeN(sortResults...)
}

func createNetworkPipeline(filename string, fileSize, chunkCount int) <-chan int {
	pipeline.Init()
	chunkSize := fileSize / chunkCount
	sortResults := []<-chan int{}
	sortAddr := []string{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)

		addr := ":" + strconv.Itoa(7000+i)
		pipeline.NetworkSink(addr, pipeline.InMemSort(source))
		sortAddr = append(sortAddr, addr)
	}

	for _, addr := range sortAddr {
		sortResults = append(sortResults, pipeline.NetworkSource(addr))
	}

	return pipeline.MergeN(sortResults...)
}
