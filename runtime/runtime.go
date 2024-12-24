package main

import (
	"fmt"
	"runtime"
	"time"
)

func main()  {
	var x int
	threads :=  runtime.GOMAXPROCS(0)

	fmt.Println("thread=", threads)
	for i:=0; i< threads; i++{
		go func() {
			defer func() {
				fmt.Println("dd-",x)
			}()
			for{x++}
		}()
	}

	time.Sleep(5*time.Second)
	fmt.Println("x=", x)
}


