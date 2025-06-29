package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/buildkite/shellwords"
	"github.com/go-cmd/cmd"
)

type Process struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`
	Command    string `json:"command"`
	Status     int64  `json:"status"`
	OrderId    int64  `json:"order_id"`

	RunStatus  string            `json:"run_status"`
	GoCmd      *cmd.Cmd          `json:"-"`
	StatusChan <-chan cmd.Status `json:"-"`
}

func (p *Process) Run(logsFun func(time int64, t string, l string)) {
	c, err := shellwords.Split(p.Command)
	if err != nil {
		fmt.Println(err)
		logsFun(time.Now().UnixMilli(), "err", err.Error())
		return
	}
	fmt.Println(c)
	p.GoCmd = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
		BeforeExec: []func(cmd *exec.Cmd){
			func(cmd *exec.Cmd) {
				SetNoWindow(cmd)
			},
		},
	}, c[0], c[1:]...)
	p.StatusChan = p.GoCmd.Start()
	go func() {
		for p.GoCmd.Stdout != nil || p.GoCmd.Stderr != nil {
			select {
			case line, open := <-p.GoCmd.Stdout:
				if !open {
					p.GoCmd.Stdout = nil
					continue
				}
				logsFun(time.Now().UnixMilli(), "out", line)
			case line, open := <-p.GoCmd.Stderr:
				if !open {
					p.GoCmd.Stderr = nil
					continue
				}
				logsFun(time.Now().UnixMilli(), "err", line)
			}
		}
	}()
}

func (p *Process) Stop() {
	if p.GoCmd != nil {
		err := p.GoCmd.Stop()
		if err != nil {
			log.Println(err)
			err := KillProcessByPID(p.GoCmd.Status().PID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
