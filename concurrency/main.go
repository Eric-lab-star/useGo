package main

import (
	"fmt"
)

func main() {
	ch := make(chan string, 1)
	timeout := make(chan bool, 1)
	go func() {
		// time.Sleep(1 * time.Second)
		timeout <- true
	}()
	go func() {
		ch <- "Hi"
	}()
	select {
	case <-ch:
		fmt.Println("channel")
		// a read from ch has occurred
	case <-timeout:
		fmt.Println("timeout")
		// the read from ch has timed out
	}
}

// may cause race condition if c.DoQuery is ready but ch is not ready
func Query(conns []Conn, query string) Result {
	ch := make(chan Result)
	for _, conn := range conns {
		go func(c Conn) {
			select {
			case ch <- c.DoQuery(query):
			default:
			}
		}(conn)
	}
	return <-ch
}
