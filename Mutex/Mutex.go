package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"yxzq.com/golib/log"
)

const (
	maxGoroutine 	= 5		//协程数
	maxDataCount	= 10	//每个协程写入的数据块个数
	maxData			= 10	//每个数据块的数据个数
	)

var protecting uint64

func init()  {
	flag.Uint64Var(&protecting, "p",1,"if use mutex")
}

func main() {
	flag.Parse()

	var buff bytes.Buffer	//数据缓冲区
	var mutex sync.Mutex

	ch := make(chan struct{}, maxGoroutine)
	for i:= 1;i<=maxGoroutine ; i++ {
		go func(id int, writer io.Writer) {
			defer func() {
				ch<- struct{}{}
			}()

			for j:=1; j<= maxDataCount; j++{
				head := fmt.Sprintf("\n[id:%d iteration:%d]", id, j)
				data := fmt.Sprintf(" %d", id*j)

				if protecting > 0{
					mutex.Lock()
				}

				_,err := writer.Write([]byte(head))
				if err != nil{
					log.Error("Write err:%v", err)
				}
				for k:=0; k<maxData ; k++  {
					writer.Write([]byte(data))
				}

				if protecting > 0{
					mutex.Unlock()
				}
			}
		}(i, &buff)
	}

	for k:= 0; k<maxGoroutine ; k++ {
		<-ch
	}
	content,_ := ioutil.ReadAll(&buff)

	fmt.Printf("protecting:%d , the content:%s", protecting, string(content))

}

/*func main()  {
	n := 0
	m := 0
	defer func() {
		fmt.Println(m,n)
	}()
	var mutex sync.Mutex
	ch := make (chan bool);

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	go func() {
		for{
			select {
			case <-ch:
				fmt.Println("ch return")
				return
			default:
				mutex.Lock()
				time.Sleep(100* time.Millisecond)
				m++
				mutex.Unlock()
			}
		}
	}()

	for i:=0;i<10;i++{
		time.Sleep(100*time.Millisecond);
		mutex.Lock()
		n++
		mutex.Unlock()
	}
	ch<-true
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
} */
