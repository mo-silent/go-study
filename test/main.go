// package main

// // import "C"

// func main() {
// 	var ch chan struct{}
// 	<-ch
// }
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// var wg sync.WaitGroup
	// foo := make(chan int)
	// bar := make(chan int)
	// closing := make(chan struct{})
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	select {
	// 	case foo <- <-bar:
	// 		println("bar")
	// 	case <-closing:
	// 		println("closing")
	// 	}
	// }()
	// // bar <- 123
	// close(closing)
	// close(bar)
	// <-foo
	// close(foo)
	// wg.Wait()
	var ar = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
	// var a []byte
	a := ar[2:5]
	a[0] = 'd'
	fmt.Println(ar)
	defer fmt.Println(1)
	defer fmt.Println(2)
	defer fmt.Println(3)
	// type Human struct {
	// 	name   string
	// 	age    int
	// 	weight int
	// }

	// type Student struct {
	// 	Human      // 匿名字段，那么默认Student就包含了Human的所有字段，
	// 	speciality string
	// 	int        // 内置类型作为匿名字段，变量名就是 int
	// }
	// jane := Student{Human: Human{"Jane", 35, 100}, speciality: "Biology", int: 1}
	// fmt.Println("Her preferred number is", jane.int)

	// Bob := Human{"Bob", 39, "sssfgf"}
	// fmt.Println("This Human is : ", Bob)

	// A interview questions
	m := make(map[int]int, 10)
	for i := 1; i <= 10; i++ {
		m[i] = i
	}
	// 闭包使用外部变量，输出会有问题
	// for k, v := range m {
	// 	go func() {
	// 		fmt.Println("k ->", k, "v ->", v)
	// 	}()
	// }
	// 解决，不要闭包直接使用外部变量，通过传参就能解决了
	for k, v := range m {
		go func(a, b int) {
			fmt.Println("k ->", a, "v ->", b)
		}(k, v)
	}

	doc := []interface{}{
		"test",
		"name",
	}
	fmt.Printf("%T\n", doc)
}

type Human struct {
	name  string
	age   int
	phone string
}

// 通过这个方法 Human 实现了 fmt.Stringer
func (h Human) String() string {
	return "❰" + h.name + " - " + strconv.Itoa(h.age) + " years -  ✆ " + h.phone + "❱"
}
