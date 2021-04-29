package goroutine_pool

import (
	"yxzq.com/golib/log"
	"sync"
	"time"

)

/*	author: xyf
*	协程池：多个任务同时并发执行，支持异步,多协程处理
*	1、初始化一个指定大小的协程池
*	2、把所有需要执行的任务放入任务列表中, 放入后会立即执行
*	3、使用完一定要关闭，不然会协程泄露
 */
//const MaxTaskCount = 500000		//最多可执行的任务个数
//const MaxChanCount = 50			//chan缓冲长度
const waitingTime 	= 500

type Task func() error

type PoolEx struct {
	queue  				chan  Task		//任务队列
	capacity 			int				//协程池容量  同时并发的最大协程个数
	chanCount			int				//chan缓冲长度
	wg 					sync.WaitGroup	//用来阻塞主协程
	mutexQueue			sync.Mutex		//队列互斥锁
	mutexClose			sync.Mutex		//协程池关闭互斥锁
	bClose				bool			//协程池是否关闭
}


//初始化一个协程池，capacity:同时允许最大并发协程数
func (p *PoolEx)Init(capacity int) {
	p.capacity 	= capacity
	p.chanCount = MaxChanCount
	p.queue 	= make(chan Task, MaxChanCount)
	p.bClose 	= false

	p.Start()
}

//启动对象池任务
func (p *PoolEx)Start()  {
	p.wg.Add(p.capacity)
	for i:= 0; i<p.capacity; i++ {
		go func() {
			defer p.wg.Done()
			for {
				p.mutexQueue.Lock()
				queueLen := len(p.queue)
				p.mutexQueue.Unlock()
				if queueLen == 0{
					p.mutexClose.Lock()
					bClose := p.bClose
					p.mutexClose.Unlock()
					if bClose{
						break
					}else{
						time.Sleep(waitingTime * time.Millisecond)
						continue
					}
				}else{
					p.mutexQueue.Lock()
					task,ok := <- p.queue		//取任务
					p.mutexQueue.Unlock()
					if !ok {
						log.Debug("task queue is closed")
						break
					}
					if err:= task(); err !=nil{			//执行任务
						log.Error("task execute Failed, Err:%+v",err)
					}
				}
			}
		}()
	}

}

func (p *PoolEx)AddTask(f Task) {
	p.mutexClose.Lock()
	isClose := p.bClose
	p.mutexClose.Unlock()

	if isClose{
		log.Error("the goroutine pool is closed")
		return
	}

	p.mutexQueue.Lock()

	if len(p.queue) == cap(p.queue){		//chan缓冲区满了需要扩容
		log.Debug("the goroutine queue len is equal capacity, so expansion")
		close(p.queue)
		tmpQue := make(chan Task, p.chanCount)
		for v := range p.queue{
			tmpQue <- v
		}
		p.queue = make(chan Task, 2*p.chanCount)		//每次缓冲区扩容两倍
		p.chanCount = 2*p.chanCount
		close(tmpQue)
		for newV := range tmpQue{
			p.queue <- newV
		}
	}

	p.queue <- f
	p.mutexQueue.Unlock()
}


func (p * PoolEx)FinishPool(){		//关闭协程池   所有任务添加完成之后一定要关闭，不然会协程泄露
	p.mutexClose.Lock()
	p.bClose = true
	p.mutexClose.Unlock()

	p.mutexQueue.Lock()
	close(p.queue)
	p.mutexQueue.Unlock()

	log.Debug("the goroutine pool Finish closed")
	p.wg.Wait()
	log.Debug("all task is done!")
}