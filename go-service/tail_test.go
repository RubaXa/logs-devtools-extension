package main

import (
	"testing"
)

func TestStart(t *testing.T) {
	m := &MockLog{}
	m.Gen(5)

	tail, lines, err := StartTail(testLogFilename, 2)
	defer tail.Close()

	if err != nil {
		t.Fatal(err)
	}

	if l := len(lines); l != 2 {
		t.Errorf("len(lines): %d != 2", l)
	}

	if v := lines[1]; v != "L4" {
		t.Error("Failed item: " + v)
	}

	if v := lines[0]; v != "L3" {
		t.Error("Failed item: " + v)
	}
}

func TestSubscribe(t *testing.T) {
	m := &MockLog{}
	m.Gen(5)

	tail, lines, err := StartTail(testLogFilename, 0)
	defer tail.Close()

	if err != nil {
		t.Fatal(err)
	}

	if l := len(lines); l != 0 {
		t.Errorf("len(lines): %d != 0", l)
	}

	s := tail.Subscribe()
	end := make(chan bool)
	newLines := make([]string, 0)

	go func() {
		for line := range s.Flow {
			newLines = append(newLines, line)
			if len(newLines) == 2 {
				s.Unsubsribe()
			}
		}

		end <- true
	}()

	m.Add("foo")
	m.Add("bar")
	<-end

	if l := len(newLines); l != 2 {
		t.Errorf("len(newLines): %d != 2", l)
	}

	if v := newLines[0]; v != "foo" {
		t.Errorf("Failed #0 -> %s\n", v)
	}

	if v := newLines[1]; v != "bar" {
		t.Errorf("Failed #1 -> %s\n", v)
	}
}

func TestMany(t *testing.T) {
	m := &MockLog{}
	m.Gen(5)

	tail1, _, err := StartTail(testLogFilename, 0)
	if err != nil {
		t.Fatal(err)
	}
	tail1.Close()
	m.Add("x")

	tail2, _, err := StartTail(testLogFilename, 0)
	if err != nil {
		t.Fatal(err)
	}

	tail2.Close()
	m.Add("y")
}
