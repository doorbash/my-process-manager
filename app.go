package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	githubUrl      string
	dbName         string
	dbHandler      *DBHandler
	processHandler *ProcessHandler
}

// NewApp creates a new App application struct
func NewApp(githubUrl, dbName string) *App {
	return &App{
		githubUrl: githubUrl,
		dbName:    dbName,
	}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx
	b.dbHandler = NewDBHandler(fmt.Sprintf("./%s", b.dbName), 3*time.Second)
	err := b.dbHandler.Init(ctx)
	if err != nil {
		log.Println(err)
	}
	b.processHandler = NewProcessHandler(func(id int64, runStatus string) {
		runtime.EventsEmit(b.ctx, "run-status", id, runStatus)
	}, func(id int64, time int64, _type string, message string) {
		runtime.EventsEmit(b.ctx, "run-log", id, time, _type, message)
	})
	b.processHandler.Start(b.ctx, b.dbHandler)
	go func() {
		systray.Run(func() {
			systray.SetIcon(Icon)
			systray.SetTitle(APP_TITLE)
			systray.SetTooltip(APP_TITLE)
			mOpen := systray.AddMenuItem(fmt.Sprintf("Show %s", APP_TITLE), "Show App")
			mOpen.Click(func() { runtime.WindowShow(b.ctx) })
			mQuit := systray.AddMenuItem("Exit", "Exit app")
			mQuit.Click(func() { systray.Quit() })
			systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(b.ctx) })
			systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })
		}, func() {
			runtime.Quit(b.ctx)
		})
	}()
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
	runtime.EventsOn(
		ctx,
		"event-from-js",
		func(optionalData ...interface{}) {
			if len(optionalData) == 0 {
				return
			}
			key, ok := optionalData[0].(string)
			if ok {
				b.onEvent(key, nil)
				return
			}
			list := optionalData[0].([]interface{})
			if len(list) == 0 {
				return
			}
			key, ok = list[0].(string)
			if !ok {
				return
			}
			if len(list) == 1 {
				b.onEvent(key, nil)
				return
			}
			b.onEvent(key, list[1])
		})
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	runtime.EventsOff(ctx, "event-from-js")
	// err := KillProcessByName(b.clientFileName) // Kill all child processes by their name!
	b.processHandler.Stop()
}

func (b *App) onEvent(key string, value interface{}) {
	fmt.Println("New Event!! key:", key, "value:", value)
}

func (b *App) OpenGithub() {
	OpenBrowser(b.githubUrl)
}

func (b *App) InsertProcess(p *Process) error {
	id, err := b.dbHandler.InsertProcess(b.ctx, p)
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.New("id is 0")
	}
	p.Id = id
	if p.Status == 1 {
		b.processHandler.AddProcess(p)
	}
	return nil
}

func (b *App) UpdateProcess(p *Process) error {
	op, err := b.dbHandler.GetProcess(b.ctx, p.Id)
	if err != nil {
		return err
	}
	err = b.dbHandler.UpdateProcess(b.ctx, p)
	if err != nil {
		return err
	}
	if op.Command != p.Command {
		b.processHandler.RemoveProcess(p.Id)
		p, err = b.dbHandler.GetProcess(b.ctx, p.Id)
		if err != nil {
			return err
		}
		if p.Status == 1 {
			b.processHandler.AddProcess(p)
		}
	}
	return nil
}

func (b *App) DeleteProcess(id int64) error {
	b.processHandler.RemoveProcess(id)
	b.processHandler.DeleteLogs(id)
	return b.dbHandler.DeleteProcess(b.ctx, id)
}

func (b *App) GetProcesses() []Process {
	list, err := b.dbHandler.GetProcesses(b.ctx, false)
	if err != nil {
		return []Process{}
	}
	plist := b.processHandler.processList.GetAll()
	// log.Println("plist:", plist)

	for i := range list {
		for j := range plist {
			if plist[j] != nil && list[i].Id == plist[j].Id {
				list[i].RunStatus = plist[j].RunStatus
				break
			}
		}
	}

	return list
}

func (b *App) RunProcess(id int64) error {
	fmt.Println("run called mother fucker")
	p, err := b.dbHandler.GetProcess(b.ctx, id)
	if err != nil {
		return err
	}
	p.Status = 1
	err = b.dbHandler.UpdateProcess(b.ctx, p)
	if err != nil {
		return err
	}
	b.processHandler.AddProcess(p)
	return nil
}

func (b *App) StopProcess(id int64) error {
	p, err := b.dbHandler.GetProcess(b.ctx, id)
	if err != nil {
		return err
	}
	p.Status = 0
	err = b.dbHandler.UpdateProcess(b.ctx, p)
	if err != nil {
		return err
	}
	b.processHandler.RemoveProcess(p.Id)
	return nil
}

func (b *App) GetLogs(id int64) []*Log {
	return b.processHandler.GetLogs(id)
}

func (b *App) ProcessesReorder(ids []int64) bool {
	return b.dbHandler.UpdateProcessesOrderId(b.ctx, ids) == nil
}

func (b *App) DeleteLogs(id int64) {
	b.processHandler.ClearLogs(id)
	runtime.EventsEmit(b.ctx, "run-log", id, "", "clear", "")
}
