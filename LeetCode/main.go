/*
Given a sorted array of distinct integers and a target value, return the index if the target is found. If not, return the index where it would be if it were inserted in order.
You must write an algorithm with O(log n) runtime complexity.

*/

package main

import "fmt"

func main() {
	input := []int{1,4,5,6}
	target:= 5
	out := searchInsert(input, target)
	fmt.Println(out)

}


func searchInsert(nums []int, target int) int {
    
}
