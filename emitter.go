package govents

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

// Initializes an Emitter to hold the registered listeners.
func InitEmitter() *Emitter {
	return &Emitter{
		Listeners: make(map[string][]*Listener),
	}
}

// Simple alias for AddListener.
func (em *Emitter) On(listener *Listener) {
	em.AddListener(listener)
}

// Appends the listener to the end of the Listeners array for the given event ID/name. No checks are made to prevent
// duplicate listeners from being added, this will cause the function to be called twice.
func (em *Emitter) AddListener(listener *Listener) {
	eventName := listener.EventName
	em.Lock()
	defer em.Unlock()
	em.Listeners[eventName] = append(em.Listeners[eventName], listener)
}

// Simple alias for RemoveListener.
func (em *Emitter) Off(listener *Listener) {
	em.RemoveListener(listener)
}

// Removes the listener from the Listeners array by matching the pointer for the given event ID/name. If multiple
// listeners are registered using the same Listener, then they will all be removed.
func (em *Emitter) RemoveListener(listener *Listener) {
	eventName := listener.EventName
	origListenerPtr := reflect.ValueOf(listener).Pointer()
	for i, curListener := range em.Listeners[eventName] {
		curListenerPtr := reflect.ValueOf(curListener).Pointer()
		if origListenerPtr != curListenerPtr {
			continue
		}
		copy(em.Listeners[eventName][i:], em.Listeners[eventName][(i+1):])
		em.Listeners[eventName][len(em.Listeners[eventName])-1] = nil
		em.Listeners[eventName] = em.Listeners[eventName][:len(em.Listeners[eventName])-1]
	}
	if len(em.Listeners[eventName]) == 0 {
		em.RemoveAllListeners(eventName)
	}
}

// Removes all listeners from the Listeners array for the given event ID/name.
func (em *Emitter) RemoveAllListeners(eventName string) {
	delete(em.Listeners, eventName)
}

// Counts the number of listeners present in the Listeners array for the given event ID/name.
func (em *Emitter) ListenerCount(eventName string) int {
	return len(em.Listeners[eventName])
}

// Lists the event ID/names which currently have listeners registered in the Listeners array.
func (em *Emitter) EventNames() []string {
	eventNames := make([]string, len(em.Listeners))
	i := 0
	for eventName := range em.Listeners {
		eventNames[i] = eventName
		i++
	}
	return eventNames
}

// Accepts an Event and fires the registered Listener function for the given event ID/name.
func (em *Emitter) Emit(event *Event) {
	for _, listener := range em.Listeners[event.EventName] {
		go listener.ListenerFn(event)
	}
}
