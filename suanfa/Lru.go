package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type LRU struct {
	maxByte int64                    //最多存储字节数，超过便会触发数据淘汰
	curByte int64                    //当前存储字节数
	list    *list.List               //go语言内置双链表存储节点
	cache   map[string]*list.Element //通过节点的key快速定位到属于哪个节点，不需要遍历双链表
	mutx    sync.RWMutex             //读写锁，保证并发安全
}

type entry struct {
	Key       string //	节点唯一标识，同时作为key存储到lru的cache里
	Value     []byte // 携带数据
	TimeStamp int64  //时间戳
}

func NewCache(maxByte int64) *LRU {
	return &LRU{
		maxByte: maxByte,
		curByte: 0,
		list:    list.New(),
		cache:   make(map[string]*list.Element),
	}
}

func (l *LRU) Get(key string) ([]byte, bool) {
	l.mutx.RLock()
	defer l.mutx.RUnlock()

	if ele, exist := l.cache[key]; exist {
		l.list.MoveToFront(ele)

		if entry, ok := ele.Value.(*entry); ok {
			return entry.Value, true
		}
	}

	return nil, false
}

func (l *LRU) Set(key string, data []byte) {
	l.mutx.Lock()
	defer l.mutx.Unlock()

	if elem, ok := l.cache[key]; ok {
		l.curByte = l.curByte - int64(len(elem.Value.(*entry).Value))
		elem.Value.(*entry).Value = data
		l.curByte += int64(len(data))
		l.list.MoveToFront(elem)
		return
	}

	l.curByte += int64(len(data))
	l.list.PushFront(&entry{Key: key, Value: data, TimeStamp: time.Now().UnixMilli()})
	l.cache[key] = l.list.Front()

	if l.curByte > l.maxByte {
		valBack := l.list.Back()
		if ent, ok := valBack.Value.(*entry); ok {
			l.curByte = l.curByte - int64(len(ent.Value))
			delete(l.cache, ent.Key)
			fomatTime := time.Unix(ent.TimeStamp, 0).Format("2006-01-02 15:04:05")
			fmt.Printf("eliminate data, key:[%s], timestamp:[%s]\n", ent.Key, fomatTime)
		}
		l.list.Remove(valBack)
	}
}

func (l *LRU) PrintLength() {
	l.mutx.RLock()
	defer l.mutx.RUnlock()

	ele := l.list.Front()
	for {
		if ele == nil {
			break
		}

		if ent, ok := ele.Value.(*entry); ok {
			fmt.Printf("Key:[%s] data:[%s]\n", ent.Key, ent.Value)
		} else {
			fmt.Printf("get value err!")
		}
		ele = ele.Next()
	}
	fmt.Printf("curByte:[%d], maxByte:[%d],Element count:[%d]\n", l.curByte, l.maxByte, len(l.cache))
}

func main() {
	c := NewCache(100)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				c.Set(fmt.Sprintf("key%d", id), []byte(fmt.Sprintf("data,data,data,data,data,data:%d", id)))
			}
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	c.PrintLength()

	/*cache := NewCache(100)
	cache.PrintLength()
	fmt.Println("\n")

	n := 0
	for {
		if n > 10 {
			break
		}
		value := []byte(fmt.Sprintf("data,data,data,data,data,data:%d", n))
		cache.Set(fmt.Sprintf("key%d", n), value)
		n++
		time.Sleep(1 * time.Second)
	}
	if data, ok := cache.Get("key9"); ok {
		fmt.Printf("key9  data:%s", string(data))
	}
	fmt.Println("\n")
	cache.PrintLength()*/
}
