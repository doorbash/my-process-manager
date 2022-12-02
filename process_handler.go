package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-cmd/cmd"
)

const (
	DELAY_BEFORE_RUN_AGAIN = 5 * time.Second
	RUN_CONFIRM_DELAY      = 3 * time.Second
)

type ProcessRunStatusFunc func(id int64, runStatus string)
type LogsFunc func(id int64, time int64, _type string, message string)

type ProcessHandler struct {
	processList   *ProcessList
	runStatusFunc ProcessRunStatusFunc
	logsFunc      LogsFunc
	ShutDown      bool
	logs          map[int64]*LogsList
	logsLock      *sync.RWMutex
}

func (ph *ProcessHandler) AddProcess(p *Process) {
	ph.RemoveProcess(p.Id)
	if !ph.processList.Add(p) {
		return
	}

	ph.logsLock.Lock()
	if ph.logs[p.Id] == nil {
		ph.logs[p.Id] = NewLogsList()
	}
	ph.logsLock.Unlock()

	go func() {
		for {
			lf := func(time int64, _type, message string) {
				ph.logsLock.RLock()
				ph.logs[p.Id].Add(&Log{
					Time:    time,
					Type:    _type,
					Message: message,
				})
				ph.logsLock.RUnlock()
				ph.logsFunc(p.Id, time, _type, message)
			}
			p.Run(lf)
			var status cmd.Status
			select {
			case <-time.After(RUN_CONFIRM_DELAY):
				p.RunStatus = "running"
				ph.runStatusFunc(p.Id, "running")
				status = <-p.StatusChan
			case status = <-p.StatusChan:
			}
			log.Println("**********************************************************************************")
			log.Println(status)
			if status.Error != nil {
				lf(time.Now().UnixMilli(), "err", status.Error.Error())
			}
			p.RunStatus = "idle"
			ph.runStatusFunc(p.Id, "idle")
			time.Sleep(DELAY_BEFORE_RUN_AGAIN)
			if p.Status == 0 {
				log.Println("ProcessHandler: AddProcess: breaking out of loop since p.Status is 0")
				break
			}
			if ph.ShutDown {
				log.Println("ProcessHandler: AddPrxoy: breaking out of loop since shutdown!")
				break
			}
		}
	}()
}

func (ph *ProcessHandler) RemoveProcess(id int64) {
	p := ph.processList.Remove(id)
	if p != nil {
		p.Status = 0
		p.Stop()
	}
}

func (ph *ProcessHandler) Start(ctx context.Context, dbHandler *DBHandler) {
	list, err := dbHandler.GetProcesses(ctx, true)
	if err != nil {
		return
	}
	for i := range list {
		p := list[i]
		ph.AddProcess(&p)
	}
}

func (ph *ProcessHandler) Stop() {
	ph.ShutDown = true
	list := ph.processList.GetAll()
	for i := range list {
		p := list[i]
		if p != nil {
			p.Stop()
		}
	}
}

func (ph *ProcessHandler) GetLogs(id int64) []*Log {
	ph.logsLock.RLock()
	defer ph.logsLock.RUnlock()
	if ph.logs[id] == nil {
		return []*Log{}
	}
	return ph.logs[id].GetAll()
}

func (ph *ProcessHandler) DeleteLogs(id int64) {
	ph.logsLock.Lock()
	defer ph.logsLock.Unlock()
	delete(ph.logs, id)
}

func (ph *ProcessHandler) ClearLogs(id int64) {
	ph.logsLock.RLock()
	ll := ph.logs[id]
	ph.logsLock.RUnlock()
	if ll != nil {
		ll.Clear()
	}
}

func NewProcessHandler(
	runStatusFunc ProcessRunStatusFunc,
	logsFunc LogsFunc,
) *ProcessHandler {
	return &ProcessHandler{
		processList:   NewProcessList(),
		runStatusFunc: runStatusFunc,
		logsFunc:      logsFunc,
		logs:          map[int64]*LogsList{},
		logsLock:      &sync.RWMutex{},
	}
}
