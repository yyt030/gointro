package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			wg.Add(1)
			go func() {
				data, err := ioutil.ReadFile(path)
				select {
				case c <- result{path, md5.Sum(data), err}:
				case <-done:
				}
				wg.Done()
			}()
			select {
			default:
				return nil
			case <-done:
				return errors.New("walk canceled")
			}
		})
		// Walk函数已经返回，启动goroutine关闭c
		go func() {
			wg.Wait()
			close(c)
		}()

		errc <- err
	}()
	return c, errc
}

func MD5All2(root string) (map[string][md5.Size]byte, error) {
	done := make(chan struct{})
	defer close(done)

	c, errc := sumFiles(done, root)

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}

	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

func main() {
	var dirName string
	if len(os.Args) >= 2 {
		dirName = os.Args[1]
	} else {
		dirName = "./"
	}

	bytes, err := MD5All2(dirName)
	if err != nil{
		panic(err)
	}
	for name,v  := range bytes{
		fmt.Printf("%x %v\n", v, name)
	}
}
