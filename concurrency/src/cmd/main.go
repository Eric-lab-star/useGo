package main

import "fmt"

func main() {
	hotLine := new(pipe)
	fmt.Println(hotLine)
}

type pipe struct {
	name    string
	channel chan string
}
