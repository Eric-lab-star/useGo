package main

import (
	"os"

	"runtime/pprof"
)

func main() {

	defer pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

type ball struct {
	count int
}

func (b *ball) counter() {
	b.count++
}
