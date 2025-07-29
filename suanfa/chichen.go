package main

import (
	"fmt"
)

/*给定三种盒装鸡块规格：7、13、29 元/盒，各盒只计“块数”不计“重量”。
  老板给你 M 元 且必须 恰好买 n 块（不能多也不能少）

即，问 ：是否存在非负整数 (a,b,c) 使得
  7a + 13b + 29c = M
  a + b + c = n
成立？

解题思路：
  两个方程联立消元：
	a = n - b - c
  	7(n - b - c) + 13b + 29c = M
  	=> 6b + 22c = M - 7n      (记 s = M - 7n)

 因此只需在 b,c ≥ 0 且 b + c ≤ n 的前提下，判断
  6b + 22c == s
  即可。
  时间复杂度 O(n) 以内即可搞定（n 一般不会特别大，常数也很小）。

*/

// Chicken 返回 (a,b,c) 和是否可行
func Chicken(M, n int) (a, b, c int, ok bool) {
	s := M - 7*n
	if s < 0 {
		return 0, 0, 0, false
	}

	// 6b + 22c = s
	for c := 0; 22*c <= s && c <= n; c++ {
		rem := s - 22*c
		if rem%6 != 0 {
			continue
		}
		b := rem / 6
		a := n - b - c
		if a >= 0 && b >= 0 {
			return a, b, c, true
		}
	}
	return 0, 0, 0, false
}

func main() {
	tests := []struct{ M, n int }{
		{50, 6},  // 可行: 7*3 + 13*2 + 29*1 = 50, 3+2+1=6
		{100, 4}, // 不可行
		{29, 1},  // 可行: 0 0 1
		{7, 1},   // 可行: 1 0 0
		{22, 2},  // 不可行
		{42, 2},
		{49, 3},
	}

	for _, t := range tests {
		a, b, c, ok := Chicken(t.M, t.n)
		if ok {
			fmt.Printf("M=%d n=%d -> %d盒7块 + %d盒13块 + %d盒29块\n",
				t.M, t.n, a, b, c)
		} else {
			fmt.Printf("M=%d n=%d -> 不可行\n", t.M, t.n)
		}
	}
}
