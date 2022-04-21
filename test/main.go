// package main

// // import "C"

// func main() {
// 	var ch chan struct{}
// 	<-ch
// }
package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup
	foo := make(chan int)
	bar := make(chan int)
	closing := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case foo <- <-bar:
			println("bar")
		case <-closing:
			println("closing")
		}
	}()
	// bar <- 123
	close(closing)
	close(bar)
	<-foo
	close(foo)
	wg.Wait()
}
