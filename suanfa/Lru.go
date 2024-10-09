package main
import (
	"fmt"
	"sync"
)


type LRU struct{
	maxByte  int64						//最多存储字节数，超过便会触发数据淘汰
	curByte  int64						//当前存储字节数
	list 	 *list.list					//go语言内置双链表存储节点
	cache     map[string]*list.Element  //通过节点的key快速定位到属于哪个节点，不需要遍历双链表
	mu 		  sync.RWMutex				//读写锁，保证并发安全
}

type Entry struct{
	Key			string		//	节点唯一标识，同时作为key存储到lru的cache里
	Value 		[]byte		// 携带数据
}


func NewCache(maxByte int64)*LRU{
	return &LRU{
		maxByte:	maxByte,
		curByte:	0,
		list:		list.New(),
		cache:		make(map[string]*list.Element),
	}
}

func (l *LRU)Get(key string)([]byte, bool){
	l.mu.RLock()
	defer l.mu.RUnLock()
	if ele,exist := l.cache[key]; exist {
		l.list.MoveToFront(ele)

		if entry,ok := ele.Value.(*Entry); ok{
			return entry.Value, true
		}
	}

	return nil,false
}


