package QpsLimit

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type QpsLimit struct {
	QpsLimit 		int32				//上限次数
	ReqCurent		int32				//当前次数
	Interval  		time.Duration	//单位间隔时间
}

func NewQPS(intval time.Duration, count int32) *QpsLimit {
	reqLimit := &QpsLimit{
		QpsLimit:   count,
		ReqCurent:  0,
		Interval: intval,
	}

	go func() {
		ticker := time.NewTicker(reqLimit.Interval)
		defer  ticker.Stop()
		for {
			<-ticker.C
			atomic.SwapInt32(&reqLimit.ReqCurent,0)
		}
	}()

	return reqLimit
}


func (q *QpsLimit)Avalib() bool{
	return !atomic.CompareAndSwapInt32(&q.ReqCurent, q.QpsLimit, q.QpsLimit)
}

func (q *QpsLimit)Increasable()  {
	atomic.AddInt32(&q.ReqCurent,	1)
}

func (q *QpsLimit)Check()bool{
	fmt.Printf("Check hit at :[%s]\n", time.Now().Format("2006-01-02 15:04:05.000"))
	if q.Avalib(){
		q.Increasable()
		return true
	}
	fmt.Printf("QPS[%d] over limit[%d]\n", q.ReqCurent, q.QpsLimit)
	return false
}

type requestLimitServers map[string]*QpsLimit

var obj *requestLimitServers = &requestLimitServers{}

func Init() * requestLimitServers {
	return  obj
}

var		qpsMutex  sync.Mutex

func Check(key string, limit int32, interval time.Duration) bool {
	q := Init()

	qpsMutex.Lock()
	defer qpsMutex.Unlock()

	if _,ok:= (*q)[key];!ok{
		(*q)[key] = NewQPS(interval, limit)
	}


	if (*q)[key].Avalib(){
		(*q)[key].Increasable()
		return true
	}

	fmt.Printf("server:%s check overlop count:%d", key, (*q)[key].ReqCurent)
	return false
}

