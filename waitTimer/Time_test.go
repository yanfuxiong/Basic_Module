package waitTimer

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimer(t * testing.T)  {
	tk := time.Tick(1*time.Second)

	i := 0
	for{
		<-tk
		println(i)
		i++
		if i > 10{
			break
		}
	}
}

func TestAfterFun(t * testing.T)  {
	c := make(chan struct{})
	tk := time.AfterFunc(5*time.Second, func() {
		fmt.Println("timer trigger")
		close(c)
	})

	if tk.C == nil{			//AfterFunc 返回的 Timer里的管道是nil
		fmt.Printf(" time.timer is nil")
	}

	//for {
		select {
		case <- tk.C:	//	会死锁
			fmt.Printf("ddd")
		case <-c:
			fmt.Printf("after fun trriger and next deal this\n")
		}
	//}

	//tmptk, ok := <-tk.C
	fmt.Println("select")

}

func TestWaitTimer(t * testing.T)  {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//time.Sleep(3 * time.Second)
		fmt.Printf("goroutine is over!\n")
	}()

	WaitTimer(&wg,time.Second)
}

func TestShowPeople(t * testing.T) {
	tch := Teacher{}
	tch.ShowA()
	tch.People.ShowA()

	//tch.ShowB()
	//tch.People.ShowB()
}

func TestPanic(t * testing.T)  {
	intChan 	:= make(chan int,1)
	stringChan 	:= make(chan string,1)
	intChan<- 1;
	stringChan <- "hello"
	select {
	case value:= <-intChan:
		fmt.Println(value)
	case value := <- stringChan:
		panic(value)
	}

}

func TestSlince(t * testing.T)  {
	s := make([]int ,5)
	s = append(s,1,2,3)
	fmt.Println(s)
}


func TestArray(t *testing.T) {
	students := parse_student()
	fmt.Printf("===================\n")
	for k,v := range students{
		fmt.Printf("key=%s,value=%+v\n", k,v)
	}
}





