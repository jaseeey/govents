package events

import (
	"reflect"
	"sync"
)

type Event struct {
	EventName string
	Args      []interface{}
}

type Listener struct {
	EventName  string
	ListenerFn func(*Event)
}

type Emitter struct {
	sync.Mutex
	Listeners map[string][]*Listener
}

func New() *Emitter {
	return &Emitter{
		Listeners: make(map[string][]*Listener),
	}
}

func (em *Emitter) On(listener *Listener) {
	em.AddListener(listener)
}

func (em *Emitter) AddListener(listener *Listener) {
	eventName := listener.EventName
	em.Lock()
	defer em.Unlock()
	em.Listeners[eventName] = append(em.Listeners[eventName], listener)
}

func (em *Emitter) Off(listener *Listener) {
	em.RemoveListener(listener)
}

func (em *Emitter) RemoveListener(listener *Listener) {
	eventName := listener.EventName
	origListenerPtr := reflect.ValueOf(listener).Pointer()
	for i, curListener := range em.Listeners[eventName] {
		curListenerPtr := reflect.ValueOf(curListener).Pointer()
		if origListenerPtr == curListenerPtr {
			copy(em.Listeners[eventName][i:], em.Listeners[eventName][(i+1):])
			em.Listeners[eventName][len(em.Listeners[eventName])-1] = nil
			em.Listeners[eventName] = em.Listeners[eventName][:len(em.Listeners[eventName])-1]
		}
	}
	if len(em.Listeners[eventName]) == 0 {
		em.RemoveAllListeners(eventName)
	}
}

func (em *Emitter) RemoveAllListeners(eventName string) {
	delete(em.Listeners, eventName)
}

func (em *Emitter) ListenerCount(eventName string) int {
	return len(em.Listeners[eventName])
}

func (em *Emitter) EventNames() []string {
	eventNames := make([]string, len(em.Listeners))
	i := 0
	for eventName := range em.Listeners {
		eventNames[i] = eventName
		i++
	}
	return eventNames
}

func (em *Emitter) Emit(event *Event) {
	for _, listener := range em.Listeners[event.EventName] {
		go listener.ListenerFn(event)
	}
}
