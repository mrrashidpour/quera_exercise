package main

type Qutex struct {
	ch chan struct{}
}

func NewQutex() *Qutex {
	q := &Qutex{
		ch: make(chan struct{}, 1),
	}
	q.ch <- struct{}{}
	return q
}

func (q *Qutex) Lock() {
	<-q.ch
}

func (q *Qutex) Unlock() {
	select {
	case q.ch <- struct{}{}:
	default:
		panic("unlock of unlocked Qutex")
	}
}
