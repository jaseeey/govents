package events

import "reflect"

var allListeners = map[string][]*Listener{}

type Event struct {
	id   string
	data interface{}
}

type Listener struct {
	eventId    string
	callableFn func(data interface{})
}

func AttachListener(listener *Listener) {
	eventId := listener.eventId
	allListeners[eventId] = append(allListeners[eventId], listener)
}

func DetachListener(listener *Listener) {
	eventId := listener.eventId
	for i, eventListener := range allListeners[eventId] {
		listener1 := reflect.ValueOf(eventListener).Pointer()
		listener2 := reflect.ValueOf(listener).Pointer()
		if listener1 == listener2 {
			copy(allListeners[eventId][i:], allListeners[eventId][(i + 1):])
			allListeners[eventId][len(allListeners[eventId])-1] = nil
			allListeners[eventId] = allListeners[eventId][:len(allListeners[eventId])-1]
		}
	}
}

func ClearAllListeners(eventId string) {
	delete(allListeners, eventId)
}

func getAllListeners(eventId string) []*Listener {
	return allListeners[eventId]
}
