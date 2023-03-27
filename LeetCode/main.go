//Remove Element

package main

import "fmt"

func main() {
	input := []int{3, 2, 2, 3}
	k := removeElements(input, 3)
	fmt.Print(k)
}

func removeElements(nums []int, val int) int {
	return len(nums)
}
