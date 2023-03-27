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
	j := 0
	if len(haystack) < len(needle) {
		return -1
	}
	for i := 0; i <= len(haystack) && i+len(needle) <= len(haystack); i++ {
		if haystack[i] == needle[j] {
			for k := i; k < i+len(needle); k++ {
				if haystack[k] != needle[j] {
					j = 0
					break
				}
				if j == len(needle)-1 {
					return i
				}
				j++
			}
		}
	}
	return -1
}
