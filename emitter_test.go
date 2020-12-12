package events

import (
	"gotest.tools/assert"
	"reflect"
	"testing"
)

const mockEventName string = "call:event"

func TestNew(t *testing.T) {
	em := New()
	_, lockExists := reflect.TypeOf(em).MethodByName("Lock")
	_, unlockExists := reflect.TypeOf(em).MethodByName("Unlock")
	assert.DeepEqual(t, em.Listeners, make(map[string][]*Listener))
	assert.Assert(t, lockExists, true)
	assert.Assert(t, unlockExists, true)
}

func TestEmitter_AddListener(t *testing.T) {
	em := New()
	listener := Listener{
		EventName:  mockEventName,
		ListenerFn: func(event *Event) {},
	}
	em.AddListener(&listener)
	assert.Equal(t, len(em.Listeners[listener.EventName]), 1)
	assert.Equal(t, em.Listeners[listener.EventName][0], &listener)
}

func TestEmitter_RemoveListener(t *testing.T) {
	em := New()
	listener := Listener{
		EventName:  mockEventName,
		ListenerFn: func(event *Event) {},
	}
	em.Listeners[listener.EventName] = append(em.Listeners[listener.EventName], &listener)
	assert.Equal(t, len(em.Listeners[listener.EventName]), 1)
	em.RemoveListener(&listener)
	assert.Equal(t, len(em.Listeners[listener.EventName]), 0)
}

func TestEmitter_EventNames(t *testing.T) {
	em := New()
	listener := Listener{
		EventName:  mockEventName,
		ListenerFn: func(event *Event) {},
	}
	expectedEventNames := []string{mockEventName}
	em.AddListener(&listener)
	assert.DeepEqual(t, em.EventNames(), expectedEventNames)
}
