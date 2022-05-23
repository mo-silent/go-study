// Author mogd 2022-05-21
// Update mogd 2022-05-21
// Description  Calculate a+b and output the sum in standard format -- that is, the digits must be separated into groups of three by commas (unless there are less than four digits).
package main

import "fmt"

func main() {
	var a, b int32
	fmt.Scan(&a, &b)
	c := a + b
	// 1、Switch to string type first
	s := fmt.Sprintf("%d", c)
	// 2、Traverse string
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c", s[i])
		// 3、The minus sign does not need to be calculated
		if s[i] == '-' {
			continue
		}
		// 4、One comma is output for every three digits
		if (i+1)%3 == len(s)%3 && i != len(s)-1 {
			fmt.Print(",")
		}
	}
	fmt.Print("\n")

}
