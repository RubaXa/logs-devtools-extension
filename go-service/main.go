package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	filename := "./tmp.log"
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	offset := int64(0)
	readLog := func() string {
		fs, _ := file.Stat()
		size := fs.Size()

		if offset > size {
			offset = 0
		}

		bytes := make([]byte, size-offset)
		n, err := file.ReadAt(bytes, offset)

		if err != nil {
			fmt.Println("[err] read:", err)
		}

		fmt.Printf("[info] Readed: %d, Offset: %d, Size: %d\n", offset, size, n)
		offset = size

		return string(bytes)
	}

	fmt.Println(readLog())

	modified := make(chan string)
	go func() {
		for filename := range modified {
			fmt.Println(filename)
			fmt.Println(readLog())
		}
	}()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					modified <- event.Name
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
