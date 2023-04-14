package main

import (
	"fmt"
	"time"
)

type Ball struct{ hits int }

func main() {
	table := make(chan *Ball)
	go player("3", table)
	go player("1", table)
	go player("2", table)
	table <- new(Ball)
	time.Sleep(1 * time.Second)
	<-table
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
