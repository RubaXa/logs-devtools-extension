package main

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

const chunkSize = 512 // byte

type Tail struct {
	Stream
	Filename string
	offset   int64
	file     *os.File
	watcher  *fsnotify.Watcher
}

func (t *Tail) Start(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	err = watcher.Add(filename)
	if err != nil {
		return
	}

	go func() {
		t.offset = t.FileSize()
		fmt.Println("Start watcher:", t.offset)

		for event := range watcher.Events {
			if event.Op&fsnotify.Write == fsnotify.Write {
				if t.IsActive() {
					t.ReadNew()
				} else {
					t.offset = t.FileSize()
				}
			}
		}

		fmt.Println("End watcher")
	}()

	t.file = file
	t.watcher = watcher
	t.Filename = filename

	return
}

func (t *Tail) Close() {
	if !t.IsActive() {
		t.file.Close()
		t.watcher.Close()

		t.file = nil
		t.watcher = nil
	}
}

func (t *Tail) ReadNew() {
	size := t.FileSize()

	fmt.Println("ReadNew:", t.offset, size)

	chunkSize := size - t.offset
	chunk := make([]byte, chunkSize)
	_, err := t.file.ReadAt(chunk, t.offset)

	if err == nil && chunkSize > 0 {
		if chunk[chunkSize-1] == '\n' {
			chunk = chunk[:chunkSize-1]
		}

		t.NotifyAll(string(chunk))
	}

	t.offset = size
}

func (t *Tail) ReadLast(n int) (lines []string) {
	if n == 0 {
		return
	}

	remains := t.FileSize()
	buf := make([]byte, 0, chunkSize)

MAIN:
	for remains != 0 {
		size := int64(chunkSize)

		if size > remains {
			size = remains
		}

		chunk := make([]byte, size)
		remains = remains - size
		_, err := t.file.ReadAt(chunk, remains)

		if err != nil {
			break
		}

		buf = append(chunk, buf...)

		lastIdx := len(buf)
		if buf[lastIdx-1] == '\n' {
			lastIdx--
		}

		for i := lastIdx - 1; i >= 0; i-- {
			if buf[i] == '\n' {
				line := string(buf[i+1 : lastIdx])
				lines = append([]string{line}, lines...)
				lastIdx = i

				if len(lines) == n {
					break MAIN
				}
			}
		}

		buf = buf[0:lastIdx]
	}

	return lines
}

func (t *Tail) FileSize() (size int64) {
	fs, err := t.file.Stat()

	if err == nil {
		size = fs.Size()
	}

	return
}

var tailStore = make(map[string]*Tail)

func StartTail(fn string, n int) (tail *Tail, lines []string, err error) {
	tail, ok := tailStore[fn]

	if !ok {
		tail = &Tail{}
		tailStore[fn] = tail
	}

	if tail.file == nil {
		err = tail.Start(fn)
		if err != nil {
			return
		}
	}

	lines = tail.ReadLast(n)
	return
}
