// package main

// import (
// 	"fmt"
// 	"os"
// 	"runtime/pprof"
// 	"time"
// )

// func main() { //main function is goroutine

// 	table := make(chan *Ball) // unbuffered channel
// 	player("ping", table)     // execute goroutine
// 	player("pong", table)
// 	table <- new(Ball)
// 	time.Sleep(2 * time.Second)

// 	time.AfterFunc(1*time.Second, func() {
// 		close(table)
// 	})
// 	pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
// }

// type Ball struct{ hit int }

// func player(name string, table chan *Ball) {
// 	go func() {
// 		for ball := range table {
// 			ball.hit++
// 			fmt.Printf("%s hits: %d\n", name, ball.hit)
// 			time.Sleep(100 * time.Millisecond)
// 			table <- ball
// 		}
// 	}()

// }
package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	table := make(chan int)
	test(table)
	table <- 1
	table <- 2
	table <- 3

	time.Sleep(1 * time.Second)
	pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
}

func test(table chan int) {
	go func(table chan int) {
		for i := range table {
			fmt.Println(i)
		}

	}(table)

}
