//Remove Element

package main

import "fmt"

func main() {
	input := []int{3, 2, 2, 3}
	k := removeElement(input, 3)
	fmt.Println(k)
}

func removeElement(nums []int, val int) int {
	i := 0
	j := len(nums) - 1
	for i <= j {
		if nums[i] == val {
			nums[i] = 0
			nums[i], nums[j] = nums[j], nums[i]
			j--
		}
		if nums[i] != val {
			i++
		}
	}
	fmt.Println(nums, i, j)
	if j+2 == i {
		return i - 1
	}
	return i
}

// [2,3]
// [3]
// [3,2]
// []
