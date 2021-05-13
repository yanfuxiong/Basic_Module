package QpsLimit

import (
	"fmt"
	"testing"
	"time"
)

func TestQpsLimit(t * testing.T){
	//qpsLimit := NewQPS(1,20)

	//idex := 0
	intval := 10
	fmt.Printf("Start at :[%s]\n", time.Now().Format("2006-01-02 15:04:05.000"))
	for{
		time.Sleep(time.Millisecond* time.Duration(intval))
	//	if !qpsLimit.Check(){
		if !Check("key",50, time.Second ) {
			//fmt.Printf("Check false count :[%d]", qpsLimit.ReqCurent)
			fmt.Printf("Check false at:[%s]\n", time.Now().Format("2006-01-02 15:04:05.000"))
			intval = 40
			break
		}

	}
	fmt.Printf("End at :[%s]\n", time.Now().Format("2006-01-02 15:04:05.000"))
}
