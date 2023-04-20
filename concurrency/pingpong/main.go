package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() { //main function is goroutine

	table := make(chan *Ball) // unbuffered channel
	go player("ping", table)  // execute goroutine
	go player("pong", table)
	table <- new(Ball)
	time.Sleep(1 * time.Second)

	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1) //print gorountine
}

type Ball struct{ hit int }

func player(name string, table chan *Ball) {

	for {
		ball := <-table
		ball.hit++
		fmt.Printf("%s hits: %d\n", name, ball.hit)
		time.Sleep(100 * time.Millisecond)
		fmt.Println(name, "ball received")
		table <- ball
		fmt.Println(name, "ball sent")
	}

}
