package main

import (
	"fmt"
	"math/rand"
)

const (
	num = 54
)

func XiPai() {
	Pai := make([]int, num)
	for i := 0; i < num; i++ {
		fmt.Scanln(&Pai[i])
	}
	fmt.Println(Pai)

	//随机数
	//rand.NewSource(time.Now().UnixNano())
	for i := num; i > 0; i-- {
		position := rand.Intn(i)
		Pai[position], Pai[i-1] = Pai[i-1], Pai[position]
	}

	fmt.Println(Pai)
}

func main() {
	XiPai()
}
