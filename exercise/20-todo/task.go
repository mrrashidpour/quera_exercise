package qtodo

import (
	"errors"
	"time"
)

type Task interface {
	DoAction()
	GetAlarmTime() time.Time
	GetAction() func()
	GetName() string
	GetDescription() string
}

type taskStruct struct {
	name        string
	description string
	alarmTime   time.Time
	action      func()
}

//func NewTask(func(), time.Time, string, string) (*struct{}, error)

func NewTask(action func(), alarmTime time.Time, name string, description string) (*taskStruct, error) {
	if action == nil {
		return nil, errors.New("action cannot be nil")
	}
	if name == "" {
		return nil, errors.New("task name cannot be empty")
	}
	if description == "" {
		return nil, errors.New("task description cannot be empty")
	}
	if alarmTime.IsZero() || alarmTime.Before(time.Now()) {
		return nil, errors.New("alarm time cannot be zero or in the past")
	}

	return &taskStruct{
		name:        name,
		description: description,
		alarmTime:   alarmTime,
		action:      action,
	}, nil
}

func (t *taskStruct) DoAction() {
	if t.action != nil {
		t.action()
	}
}

func (t *taskStruct) GetAlarmTime() time.Time {
	return t.alarmTime
}

func (t *taskStruct) GetAction() func() {
	return t.action
}

func (t *taskStruct) GetName() string {
	return t.name
}

func (t *taskStruct) GetDescription() string {
	return t.description
}
