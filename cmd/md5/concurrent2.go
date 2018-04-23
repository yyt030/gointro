package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type result2 struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk canceld")
			}
			return nil
		})
	}()
	return paths, errc
}

func digester(done <-chan struct{}, paths <-chan string, c chan<- result2) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case c <- result2{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}

func MD5All3(done <-chan struct{}, paths <-chan string) chan result2 {
	// Start some goroutine
	c := make(chan result2)
	var wg sync.WaitGroup
	const numDigesters = 20
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(done, paths, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}

func main() {
	done := make(chan struct{})

	files, err := walkFiles(done, ".")
	if len(err) > 0 {
		panic(err)
	}

	results := MD5All3(done, files)
	for r := range results {
		fmt.Printf("%x %v\n", r.sum, r.path)
	}
}
