package main

import "fmt"

func main() {
	myMap := map[string]int{
		"apple":  1,
		"banana": 2,
		"grape":  3,
	}
	mutate(myMap)
	fmt.Println(myMap)

}

func mutate(myMap map[string]int) {
	for k, v := range myMap {
		myMap[k] = v + 1
	}
}
