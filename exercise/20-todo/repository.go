package qtodo

import (
	"errors"
	"sync"
)

type Database interface {
	GetTaskList() []Task
	GetTask(name string) (Task, error)
	SaveTask(task Task) error
	DelTask(name string) error
}

type memoryDB struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

func NewDatabase() Database {
	return &memoryDB{tasks: make(map[string]Task)}
}

func (db *memoryDB) GetTaskList() []Task {
	db.mu.RLock()
	defer db.mu.RUnlock()

	list := make([]Task, 0, len(db.tasks))
	for _, t := range db.tasks {
		list = append(list, t)
	}
	return list
}

func (db *memoryDB) GetTask(name string) (Task, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if t, ok := db.tasks[name]; ok {
		return t, nil
	}
	return nil, errors.New("task not found")
}

func (db *memoryDB) SaveTask(task Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.tasks[task.GetName()] = task
	return nil
}

func (db *memoryDB) DelTask(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.tasks[name]; !ok {
		return errors.New("task not found")
	}
	delete(db.tasks, name)
	return nil
}
