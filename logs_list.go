package main

import (
	"container/list"
	"sync"
)

const LOG_BUFFER_SIZE = 200

type LogsList struct {
	list *list.List
	lock *sync.Mutex
}

func (ll *LogsList) GetAll() []*Log {
	ll.lock.Lock()
	defer ll.lock.Unlock()
	ret := make([]*Log, 0)
	for e := ll.list.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Log))
	}
	return ret
}

func (ll *LogsList) Add(log *Log) {
	ll.lock.Lock()
	defer ll.lock.Unlock()
	for ll.list.Len() >= LOG_BUFFER_SIZE {
		ll.list.Remove(ll.list.Front())
	}
	ll.list.PushBack(log)
}

func NewLogsList() *LogsList {
	return &LogsList{
		list: list.New(),
		lock: &sync.Mutex{},
	}
}
