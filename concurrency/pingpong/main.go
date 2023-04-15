package main

import (
	"fmt"
	"time"
)

func main() {

	table := make(chan *Ball)
	go player("ping", table)
	go player("pong", table)
	table <- new(Ball)
	time.Sleep(1 * time.Second)
	panic("show goroutines")
}

type Ball struct{ hit int }

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hit++
		fmt.Printf("%s hits: %d\n", name, ball.hit)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}

}
