package main

import "fmt"

func main() {
	data := []int{1, 2, 3, 4, 5, 6}
	list1 := []int{}
	list2 := []int{}
	c1 := make(chan bool)
	c2 := make(chan bool)
	c3 := make(chan bool)
	go func() {

		go func() {
			for _, v := range data {
				list1 = append(list1, v)
			}
			c1 <- true
		}()
		go func() {
			for _, v := range data {
				list2 = append(list2, v)
			}
			c2 <- true
		}()
		if <-c1 && <-c2 {
			c3 <- true
		}
	}()
	select {
	case <-c3:
		fmt.Println(list1)
		fmt.Println(list2)
	}

}
