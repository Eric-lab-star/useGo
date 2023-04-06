package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c := fanIn(boring("jon:"), boring("ann:"))
	for i := 0; i < 10; i++ {
		msg1 := <-c
		fmt.Println(msg1.str)
		msg2 := <-c
		fmt.Println(msg2.str)
		msg1.wait <- true
		msg2.wait <- true
	}

	fmt.Println("you are boring. I am leaving.")
}

func boring(msg string) <-chan message {
	c := make(chan message)
	waitforit := make(chan bool)
	go func() {
		for i := 0; ; i++ {
			c <- message{fmt.Sprintf("%s %d", msg, i), waitforit}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-waitforit
		}
	}()
	return c
}

func fanIn(jon, ann <-chan message) <-chan message {
	c := make(chan message)
	go func() {
		for {
			c <- <-jon
		}
	}()
	go func() {
		for {
			c <- <-ann
		}
	}()
	return c
}

type message struct {
	str  string
	wait chan bool
}
