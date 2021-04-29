package goroutine_pool

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"yxzq.com/golib/log"
)

func TestNewPool(t * testing.T){
	log.SetLogFileName("./test.log")
	log.SetLogLevel(log.LOGLEVEL_DEBUG)

	taskPool := New(4)
	taskPool.AddTask(func() error {
		return testA()
	})
	taskPool.AddTask(func() error{
		return testB()
	})
	taskPool.AddTask(func() error {
		return testC()
	})
	taskPool.AddTask(func() error {
		return testD()
	})
	taskPool.AddTask(func()error{
		return testE()
	})

	i := 0
	for {
		if i==5{
			break
		}
		taskPool.AddTask(func()error {
			return testRand()
		})
		i++
	}

	taskPool.AddTask(func()error{
		return testE()
	})


	taskPool.Run("TestNewPool")

	log.Debug("[%+v]",taskPool)
}

func testA() error  {
	fmt.Println("aaaaaa")
	return nil
}
func testB() error {
	fmt.Println("bbbbbb")
	return nil
}
func testC() error {
	fmt.Println("cccccc")
	return nil
}
func testD() error {
	fmt.Println("dddddd")
	return nil
}
func testE() error  {
	fmt.Println("eeeeee")
	return nil
}
func testRand() error {
		time.Sleep(time.Duration(rand.Intn(5))*time.Second)
		nSecond := time.Now().Unix()
		fmt.Println(time.Unix(nSecond,0).Format("2006-01-02 15:03:04"))
	return nil
}

func TestPoolEx(t * testing.T)  {
	log.SetLogFileName("./test.log")
	log.SetLogLevel(log.LOGLEVEL_DEBUG)

	rand.Seed(time.Now().UnixNano())

	var pool PoolEx
	pool.Init(2)
	pool.AddTask(func() error {
		return testA()
	})
	pool.AddTask(func() error {
		return testB()
	})
	pool.AddTask(func() error {
		return testC()
	})
	pool.AddTask(func() error {
		return testD()
	})
	pool.AddTask(func() error {
		return testE()
	})

	nCount := 0
	for{
		if nCount > 5{
			break
		}
		pool.AddTask(func() error {
			return testRand()
		})
		nCount++
	}
	pool.FinishPool()
	//time.Sleep(time.Second * 5)
}