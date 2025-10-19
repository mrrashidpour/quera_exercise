package qtodo

import (
	"errors"
	"sync"
	"time"
)

type App interface {
	StartTask(name string) error
	StopTask(name string)
	AddTask(name, description string, alarmTime time.Time, action func(), temporary bool) error
	DelTask(name string) error
	GetTaskList() []Task
	GetTask(name string) (Task, error)
	GetActiveTaskList() []Task
}

type appStruct struct {
	db        Database
	active    map[string]bool
	activeMu  sync.RWMutex
	stopChans map[string]chan bool
	stopMu    sync.RWMutex
}

func NewApp(db Database) App {
	return &appStruct{
		db:        db,
		active:    make(map[string]bool),
		stopChans: make(map[string]chan bool),
	}
}

func (app *appStruct) AddTask(name, description string, alarmTime time.Time, action func(), temporary bool) error {
	task, err := NewTask(action, alarmTime, name, description)
	if err != nil {
		return err
	}
	err = app.db.SaveTask(task)
	if err != nil {
		return err
	}

	if temporary {
		go func() {
			delay := time.Until(task.GetAlarmTime())
			if delay > 0 {
				time.Sleep(delay)
			}
			task.DoAction()
			app.DelTask(name)
		}()
	}

	return nil
}

func (app *appStruct) DelTask(name string) error {
	app.StopTask(name)
	return app.db.DelTask(name)
}

func (app *appStruct) GetTaskList() []Task {
	return app.db.GetTaskList()
}

func (app *appStruct) GetTask(name string) (Task, error) {
	return app.db.GetTask(name)
}

func (app *appStruct) StartTask(name string) error {
	task, err := app.db.GetTask(name)
	if err != nil {
		return err
	}

	app.activeMu.Lock()
	if app.active[name] {
		app.activeMu.Unlock()
		return errors.New("task already started")
	}
	app.active[name] = true
	app.activeMu.Unlock()

	stopChan := make(chan bool, 1)
	app.stopMu.Lock()
	app.stopChans[name] = stopChan
	app.stopMu.Unlock()

	go func() {
		delay := time.Until(task.GetAlarmTime())
		if delay > 0 {
			select {
			case <-time.After(delay):
				app.activeMu.RLock()
				if app.active[name] {
					task.DoAction()
				}
				app.activeMu.RUnlock()
			case <-stopChan:
				return
			}
		} else {
			app.activeMu.RLock()
			if app.active[name] {
				task.DoAction()
			}
			app.activeMu.RUnlock()
		}
	}()

	return nil
}

func (app *appStruct) StopTask(name string) {
	app.activeMu.Lock()
	defer app.activeMu.Unlock()

	if !app.active[name] {
		return
	}

	app.active[name] = false

	app.stopMu.Lock()
	if ch, ok := app.stopChans[name]; ok {
		ch <- true
		close(ch)
		delete(app.stopChans, name)
	}
	app.stopMu.Unlock()
}

func (app *appStruct) GetActiveTaskList() []Task {
	app.activeMu.RLock()
	defer app.activeMu.RUnlock()

	activeTasks := []Task{}
	for name, running := range app.active {
		if running {
			task, err := app.db.GetTask(name)
			if err == nil {
				activeTasks = append(activeTasks, task)
			}
		}
	}
	return activeTasks
}
