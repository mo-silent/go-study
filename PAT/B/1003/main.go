package main

import "fmt"

func main() {
	var n int
	_, _ = fmt.Scan(&n)
	for i := 0; i < n; i++ {
		var (
			s          string
			lineBreaks = "\n" // 换行符
			numP       = 0    // P 的个数
			numT       = 0    // T 的个数
			other      = 0    // 其他字符串的个数
			locP       = -1   // P 在字符串的位置
			locT       = -1   // T 在字符串的位置
		)
		//fmt.Println(locT,locP)
		_, _ = fmt.Scan(&s)
		for j, k := range s {
			if k == 'P' {
				numP++
				locP = j
			} else if k == 'T' {
				numT++
				locT = j
			} else if k != 'A' {
				other++
			}
		}
		// 最后一个字符串的判定输出结果不换行
		if i == n-1 {
			lineBreaks = ""
		}
		// P 个数不为 1, T 的个数不为 1
		// 存在非 P A T 字符, P 和 T 之间没有字符
		if (numP != 1) || (numT != 1) || (other != 0) || (locT-locP <= 1) {
			fmt.Printf("NO" + lineBreaks)
			continue
		}
		// x, y , z 定义含义，见题目解析
		x := locP
		y := locT - locP - 1
		z := len(s) - locT - 1
		// 条件 2 成立的条件
		if z-x*(y-1) == x {
			fmt.Printf("YES" + lineBreaks)
		} else {
			fmt.Printf("NO" + lineBreaks)
		}
	}
}
