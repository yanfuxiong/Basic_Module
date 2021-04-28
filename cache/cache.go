package cache

import (
	"sync"
	"time"
)


/*	author: xyf
*   缓存管理器
*	1、根据需要初始化一个缓存管理对象
*	2、把需要缓存的数据调用Set方法存入、Get方法读取、Clear方法清空
*	3、使用完需要调用FinishGC方法关闭，不然会协程泄露
 */

type cache struct {
	mutexCount 		sync.Mutex		//缓存数量计数器互斥锁
	mp 				*sync.Map
	cacheTimeout 	time.Duration	//缓存过期时间	为0 则表示不清缓存
	gcInterval 		time.Duration	//gc间隔时间
	count 			uint64			//缓存数据当前数据量
	nMaxCount		uint64			//缓存数据最大数量， 0表示不做限制

	ch				chan struct{}	//用于退出gc 结束缓存使用
}

type cacheItem struct {
	cacheTime 	time.Time			//缓存时间戳
	value 		interface{}
}

func New(cacheTimeout, gcInterval time.Duration, max uint64)*cache{
	if gcInterval < 0{
		panic("the gcInterval must be greater than zero")
	}

	if cacheTimeout < 0{
		panic("the cacheTimeout must be greater than zero")
	}

	c := &cache{
		mutexCount:   sync.Mutex{},
		mp:           &sync.Map{},
		cacheTimeout: cacheTimeout,
		gcInterval:   gcInterval,
		count:        0,
		nMaxCount:    max,
		ch: 		  make(chan struct{}),
	}

	c.StartGC()

	return c
}

func (c *cache)StartGC()  {
	if c.cacheTimeout <= 0{
		return
	}

	go func() {
		t := time.NewTimer(c.gcInterval)
		defer t.Stop()
		for{
			select{
			case <-t.C:
				c.DoGC()
			case <-c.ch:
				return
			}
		}
	}()
}

func (c *cache)FinishGC()  {
	close(c.ch)
}

var gcDone bool = true

func (c *cache)DoGC()  {

	if c.cacheTimeout <= 0{
		return
	}

	if !gcDone {
		return
	}

	gcDone = false

	c.mp.Range(func(key, value interface{}) bool {
		item := value.(cacheItem)
		if time.Since(item.cacheTime) >= c.cacheTimeout{
			c.mp.Delete(key)

			c.mutexCount.Lock()
			c.count--
			c.mutexCount.Unlock()
		}
		return true
	})

	gcDone = true
}

func (c *cache)Set(key,value interface{}) bool {
	if ( c.nMaxCount >0 && c.count>= c.nMaxCount){		//数据量达到最大数之后则GC一次再存

		c.DoGC()
		if ( c.nMaxCount>0 && c.count>= c.nMaxCount){
			return false
		}
	}

	Item := cacheItem{
		cacheTime: time.Now(),
		value:     value,
	}
	c.mp.Store(key, Item)


	c.mutexCount.Lock()
	c.count++
	c.mutexCount.Unlock()

	return true
}

func (c *cache)Get(key interface{}) (interface{}, bool){

	value,ok := c.mp.Load(key)
	if ok{
		item := value.(cacheItem)

		return item.value, true
	}
	return nil, false
}

func (c* cache)Clear(){

	c.mp.Range(func(key, value interface{}) bool {
		c.mp.Delete(key)
		return true
	})

	c.mutexCount.Lock()
	c.count = 0
	c.mutexCount.Unlock()
}