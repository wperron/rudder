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
	f = flag.Bool("f", false, "follow mode")
)

func main() {
	flag.Parse()

	positional := flag.Args()
	if len(positional) == 0 || positional[0] == "-" {
		Tail(os.Stdin, os.Stdout, *f)
	} else {
		for _, p := range positional {
			r, err := os.Open(p)
			if err != nil {
				fmt.Printf("opening file: %s", err)
				os.Exit(1)
			}
			defer r.Close()
			fmt.Printf("==> %s <==\n", p)
			Tail(r, os.Stdout, *f)
			fmt.Print("\n")
		}
	}
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

func Tail(r io.Reader, w io.Writer, follow bool) {
	scanner := bufio.NewScanner(r)

	buf := NewCircularBuffer(*n)
	for scanner.Scan() {
		buf.Push(scanner.Text())
	}

	fmt.Fprintf(w, strings.Join(buf.ReadAll(), "\n"))
}
