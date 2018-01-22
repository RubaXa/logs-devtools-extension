package main

import (
	"fmt"
	"log"
)

func main() {
	filename := "./tmp.log"
	_, lines, err := StartTail(filename, 3)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Lines:", lines)

	// watcher, err := fsnotify.NewWatcher()

	// defer watcher.Close()

	// file, err := os.Open(filename)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer file.Close()

	// offset := int64(0)
	// readLog := func() string {
	// 	fs, _ := file.Stat()
	// 	size := fs.Size()

	// 	if offset > size {
	// 		offset = 0
	// 	}

	// 	bytes := make([]byte, size-offset)
	// 	n, err := file.ReadAt(bytes, offset)

	// 	if err != nil {
	// 		fmt.Println("[err] read:", err)
	// 	}

	// 	fmt.Printf("[info] Readed: %d, Offset: %d, Size: %d\n", offset, size, n)
	// 	offset = size

	// 	return string(bytes)
	// }

	// fmt.Println(readLog())

	// modified := make(chan string)
	// go func() {
	// 	for filename := range modified {
	// 		fmt.Println(filename)
	// 		fmt.Println(readLog())
	// 	}
	// }()

	// done := make(chan bool)

	// err = watcher.Add(filename)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// <-done
}
