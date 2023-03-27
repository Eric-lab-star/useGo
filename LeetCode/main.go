//Find the index of the first Occurence in a String

package main

import "fmt"

func main() {
	haystck := "a"
	needle := "a"
	ret := strStr(haystck, needle)
	fmt.Println(ret)

}

func strStr(haystack string, needle string) int {
	if len(haystack) < len(needle) {
		return -1
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}
