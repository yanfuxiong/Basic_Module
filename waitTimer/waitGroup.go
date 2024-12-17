package waitTimer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"time"
)

//超时的wait
func WaitTimer(wg *sync.WaitGroup, duration time.Duration) bool {

	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	tk := time.NewTimer(duration)
	defer tk.Stop()

	select {
	case <-tk.C:
			fmt.Println("timer is select")
			return true
	case <-c:
			fmt.Println("wg is return")
			return false
	}

	return false
}

type People struct { }

func (p *People)ShowA()  {
	fmt.Println("ShowA")
	p.ShowB()
}

func (p * People)ShowB()  {
	fmt.Println("ShowB")
}

type Teacher struct {
	People
}

func (t *Teacher)ShowA()  {
	fmt.Println("TeacherShowA")
}
func (t *Teacher)ShowB()  {
	fmt.Println("TeacherShowB")
}

type student struct {
	Name 	string
	Age 	int
}
func parse_student()map[string]*student {
	m := make(map[string]*student)
	stus := []student{
		{"zhangsan",23},
		{"lisi", 25},
		{"wangwu", 32},
	}

	for _,stu := range stus{
		m[stu.Name] = &stu
		fmt.Printf("%v", &stu)
	}
	return m
}

func echo(wr http.ResponseWriter, r *http.Request)  {
	msg,err := ioutil.ReadAll(r.Body)
	if err != nil{
		wr.Write([]byte("echo error"))
		return
	}

	writeLen, err := wr.Write(msg)
	if err != nil || len(msg)!= writeLen{
		fmt.Printf("write len:%d",writeLen)
		//return
	}
}

func main()  {
	/*http.HandleFunc("/", echo)
	err := http.ListenAndServe("172.16.7.17:8080",nil)
	if err != nil{
		fmt.Printf("err")
	}*/

	var wg sync.WaitGroup
	wg.Add(1)
	WaitTimer()
}