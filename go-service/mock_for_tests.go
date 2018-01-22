package main

import (
	"fmt"
	"os"
)

const testLogFilename = "./test_log.txt"

type MockLog struct {
	File     *os.File
	Filename string
}

func (m *MockLog) Open(clear bool) {
	if clear {
		os.Remove(testLogFilename)
		m.File, _ = os.Create(testLogFilename)
	} else {
		m.File, _ = os.OpenFile(testLogFilename, os.O_RDWR, 0644)
	}
}

func (m *MockLog) Write(v string) {
	m.File.WriteString(v)
}

func (m *MockLog) Writeln(v string) {
	m.Write(v + "\n")
}

func (m *MockLog) Close() {
	m.File.Close()
}

func (m *MockLog) Add(v string) {
	m.Open(false)
	m.File.Seek(0, 2)
	m.Writeln(v)
	m.Close()
}

func (m *MockLog) Gen(n int) {
	m.Open(true)
	for i := 0; i < n; i++ {
		m.Writeln(fmt.Sprintf("L%d", i))
	}
	m.Close()
}
