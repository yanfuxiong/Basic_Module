package goroutine_pool

import (
	"sync"
	"time"
	"yxzq.com/golib/log"
)

/*
*	协程池：多个任务同时并发执行
*	1、初始化一个协程池
*	2、把所有需要执行的任务放入任务列表中
*	3、启动运行，遍历任务列表，同时执行完毕退出
*/
const MaxTaskCount = 500000		//最多可执行的任务个数
const MaxChanCount = 50			//chan缓冲长度

type funTask func() error

type Pool struct {
	queue  				chan  funTask	//任务队列
	capacity 			int				//协程池容量  同时并发的最大任务个数
	chanCount			int				//chan缓冲长度
	addTotalCount		int				//所添加的总的个数
	curRunningCount 	int				//当前正在进行中的任务数
	wg 					sync.WaitGroup	//用来阻塞
	mutexQueue			sync.Mutex		//队列互斥锁
	mutexTaskCount		sync.Mutex		//任务计数器互斥锁
}


//新建一个协程任务池，capacity:同时允许最大并发数
func New(capacity int) *Pool{
	return &Pool{
		queue: make(chan funTask,	50),
		capacity:capacity,
		addTotalCount:0,
		curRunningCount:0,
		chanCount:	50,
	}
}

func (p *Pool)AddTask(f funTask) {
	if p.addTotalCount > MaxTaskCount{
		log.Error("add task to goroutine is over 500000, Please check...")
		return
	}

	p.mutexQueue.Lock()
	if len(p.queue) == cap(p.queue){		//chan缓冲区满了需要扩容
		log.Debug("queue len is equal capacity, so expansion")
		close(p.queue)
		tmpQue := make(chan funTask, p.chanCount)
		for v := range p.queue{
			tmpQue <- v
		}

		p.queue = make(chan funTask, 2*p.chanCount)		//每次缓冲区扩容两倍
		p.chanCount = 2*p.chanCount
		close(tmpQue)
		for newV := range tmpQue{
			p.queue <- newV
		}
	}

	p.queue <- f
	p.mutexQueue.Unlock()

	p.addTotalCount++
}

func (p *Pool)Run(serverName string) {
	p.CloseQueue()		//一次启动所有任务并执行 后续不能再添加任务了
	for{
		p.mutexTaskCount.Lock()
		nRunCount := p.curRunningCount
		p.mutexTaskCount.Unlock()

		if nRunCount > p.capacity{			//超过同时并发任务数量 需要暂缓一下
			log.Debug(" %s the running task is over capacity,so sleep 300ms",serverName)
			time.Sleep(time.Millisecond*300)
			continue
		}
		p.mutexQueue.Lock()
		task,ok := <- p.queue
		p.mutexQueue.Unlock()
		if !ok{
			log.Error("queue is close")
			break
		}
		p.wg.Add(1)
		p.mutexTaskCount.Lock()
		p.curRunningCount++
		p.mutexTaskCount.Unlock()
		go func() {
			defer p.wg.Done()
			err := task()
			if err != nil{
				log.Error("%s  Task exectue error:[%+v] ",serverName, err)
			}

			p.mutexTaskCount.Lock()
			p.curRunningCount--
			p.mutexTaskCount.Unlock()
		}()
	}

	p.wg.Wait()
	log.Debug("Run finished")
}

func (p * Pool)CloseQueue(){	//需要关闭
	p.mutexQueue.Lock()
	close(p.queue)
	p.mutexQueue.Unlock()
	log.Debug(" queue close ")
}