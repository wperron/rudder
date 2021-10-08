package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	n = flag.Int("n", 10, "number of lines")
)

func main() {
	flag.Parse()

	positional := flag.Args()
	var r io.Reader
	if len(positional) == 0 || positional[0] == "-" {
		r = os.Stdin
	} else {
		var err error
		// TODO(wperron) should open each file if more than one is specified
		r, err = os.Open(positional[0])
		if err != nil {
			fmt.Printf("opening file: %s", err)
			os.Exit(1)
		}
	}

	Tail(r, os.Stdout)
}

type CircularBuffer struct {
	inner []string
	length int
	pos int
}

func NewCircularBuffer(l int) *CircularBuffer {
	return &CircularBuffer{
		inner: make([]string, l, l),
		length: l,
	}
}

func (cb *CircularBuffer) Push(s string) {
	cb.inner[cb.pos] = s
	cb.pos++
	if cb.pos == cb.length {
		cb.pos = 0
	}
}

func (cb *CircularBuffer) ReadAll() []string {
	return append(cb.inner[:cb.pos], cb.inner[cb.pos:]...)
}

func Tail(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	buf := NewCircularBuffer(*n)
	for scanner.Scan() {
		buf.Push(scanner.Text())
	}

	fmt.Fprintf(w, strings.Join(buf.ReadAll(), "\n"))
}
