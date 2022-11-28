package main

import (
	"container/list"
	"sync"
)

type ProcessList struct {
	data *list.List
	lock *sync.Mutex
}

func (pl *ProcessList) GetAll() []*Process {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	ret := make([]*Process, 0)
	for e := pl.data.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Process))
	}
	return ret
}

func (pl *ProcessList) Add(p *Process) bool {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	for e := pl.data.Front(); e != nil; e = e.Next() {
		if e.Value.(*Process).Id == p.Id {
			return false
		}
	}
	pl.data.PushBack(p)
	return true
}

func (pl *ProcessList) Remove(id int64) *Process {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	for e := pl.data.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Process)
		if p.Id == id {
			pl.data.Remove(e)
			return p
		}
	}
	return nil
}

func (pl *ProcessList) Len() int {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	return pl.data.Len()
}

func (pl *ProcessList) Front() *Process {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	return pl.data.Front().Value.(*Process)
}

func (pl *ProcessList) GetById(id int64) *Process {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	for e := pl.data.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Process)
		if p.Id == id {
			return p
		}
	}
	return nil
}

func NewProcessList() *ProcessList {
	return &ProcessList{
		data: list.New(),
		lock: &sync.Mutex{},
	}
}
