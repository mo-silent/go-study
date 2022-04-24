// @Title  go-generics
// @Description  Generic usage example
// @Description  泛型使用示例
// @Author  mogd  20220423 14:44 CST
// @Update  mogd  20220423 15:41 CST
package main

import "fmt"

// Count Calculated structs. It supports string, int, int64 and float64
//
// @Description A computational struct with two variable
// @Description 一个计算结构体，有两个变量
type Count[T string | int | int64 | float64] struct {
	A T
	B T
}

// CustomizationGenerics custom type constraint
//
// @Description custom type constraint, which are type restrictions
// @Description ~is a new symbol added to Go 1.18, and the ~ indicates that the underlying type is all types of T. ~ is pronounced astilde in English
// @Description 自定义泛型，即类型限制
// @Desciption ~ 是 Go 1.18 新增的符号，~ 表示底层类型是T的所有类型。~ 的英文读作 tilde
//
// @Example With the addition of ~, MyInt can be used, otherwise there will be type mismatch
// @Example 加上 ~，那么 MyInt 自定义的类型能够被使用，否则会类型不匹配
type CustomizationGenerics interface {
	~int | ~int64
}

// MyInt Custom type
type MyInt int

// MyChan Custom generics chan type
type MyChan[T int | string] chan T

// Add sums the values of T. It supports string, int, int64 and float64
//
// @Description A simple additive generic function
// @Description 一个简单的加法泛型函数
// @parameter	a, b	T string | int | int64 | float64	"generics parameter"
// @return		c		T string | int | int64 | float64	"generics return"
func Add[T string | int | int64 | float64](a, b T) T {
	return a + b
}

// Sub sub the values of T. It supports CustomizationGenerics
//
// @Description A simple subtraction function. It is used to test custom type constraint
// @Description 一个简单的减法函数，用来测试自定义泛型
// @parameter	a, b	T CustomizationGenerics	"custom type constraint parameter"
// @return		c		T CustomizationGenerics	"custom type constraint return"
func Sub[T CustomizationGenerics](a, b T) T {
	return a - b
}

func main() {
	a := MyInt(1)

	count := Count[int]{1, 2}
	fmt.Println(Add(count.A, count.B))
	// @Interpretation if the variable definition in CustomizationGenerics add ~, a can be used. otherwise the type conflicts
	// @Interpretation 如果 CustomizationGenerics 中的变量定义加上 ~，a 能够被使用，否则类型冲突
	fmt.Println(Sub(a, 0))
}

// Official example

// package main

// import "fmt"

// // SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// // as types for map values.
// func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
// 	var s V
// 	for _, v := range m {
// 		s += v
// 	}
// 	return s
// }

// func main() {
// 	// Initialize a map for the integer values
// 	ints := map[string]int64{
// 		"first":  34,
// 		"second": 12,
// 	}

// 	// Initialize a map for the float values
// 	floats := map[string]float64{
// 		"first":  35.98,
// 		"second": 26.99,
// 	}

// 	fmt.Printf("Generic Sums: %v and %v\n",
// 		SumIntsOrFloats[string, int64](ints),
// 		SumIntsOrFloats[string, float64](floats))
// }
