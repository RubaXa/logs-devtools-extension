package main

import (
	"os"
)

type Tail struct {
	File     *os.File
	Filename string
	Offset   int64
	Lines    chan string
}

func (t *Tail) Watch(filename string, n int) (err error) {
	file, err := os.Open(filename)

	if err == nil {
		t.File = file
		t.Filename = filename
	}

	return
}

func (t *Tail) size() (size int64) {
	fs, err := t.File.Stat()

	if err == nil {
		size = fs.Size()
	}

	return
}
