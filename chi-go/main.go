package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// c := fanIn(boring("jon:"), boring("ann:"))
	// for i := 0; i < 10; i++ {
	// 	msg1 := <-c
	// 	fmt.Println(msg1.str)
	// 	msg2 := <-c
	// 	fmt.Println(msg2.str)
	// 	msg1.wait <- true
	// 	msg2.wait <- true
	// }
	quit := make(chan bool)
	c := boring("jon:", quit)

	for i := 10; i <= 0; i-- {

		fmt.Println(<-c)
	}
	quit <- true

	fmt.Println("you are boring. I am leaving.")
}

func boring(msg string, quit chan bool) <-chan message {
	c := make(chan message)
	waitforit := make(chan bool)
	go func() {
		for i := 0; ; i++ {
			select {
			case c <- message{fmt.Sprintf("%s %d", msg, i), waitforit}:
				fmt.Println(i)
			case <-quit:
				return

			}

			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			// <-waitforit
		}
	}()
	return c
}

func fanIn(jon, ann <-chan message) <-chan message {
	c := make(chan message)
	go func() {
		for {
			select {
			case j := <-jon:
				c <- j
			case a := <-ann:
				c <- a
			}
		}
	}()

	return c
}

type message struct {
	str  string
	wait chan bool
}
