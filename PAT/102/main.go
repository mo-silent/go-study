package main

import (
	"fmt"
)

func main() {
	var (
		k     int
		n     int16
		e     float32
		count int16 = 0
		a           = make([]float32, 1010)
	)
	for j := 0; j < 2; j++ {
		fmt.Scan(&k)
		for i := 0; i < k; i++ {
			fmt.Scan(&n, &e)
			if e == 0 {
				continue // The resulting index added up is 0, and the quantity is minus one
			}
			if a[n] == 0 {
				count += 1 // The first polynomial plus 1
			}
			a[n] += e // The same polynomials are added
			if a[n] == 0 {
				count -= 1 // The resulting index added up is 0, and the quantity is minus one
			}
		}
	}
	fmt.Printf("%d", count)
	for i := 1009; i >= 0; i-- {
		if a[i] != 0 {
			fmt.Printf(" %d %.1f", i, a[i])
		}
	}
	fmt.Printf("\n")
}
