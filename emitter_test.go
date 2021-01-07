package govents

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"sync"
	"testing"
)

const mockEventName string = "call:event"

func TestInitEmitter(t *testing.T) {
	em := InitEmitter()
	_, lockExists := reflect.TypeOf(em).MethodByName("Lock")
	_, unlockExists := reflect.TypeOf(em).MethodByName("Unlock")
	assert.EqualValues(t, em.Listeners, make(map[string][]*Listener))
	assert.True(t, lockExists)
	assert.True(t, unlockExists)
}

func TestAddsListenerForEvent(t *testing.T) {
	em := InitEmitter()
	listener := Listener{
		EventName:  mockEventName,
		ListenerFn: func(event *Event) {},
	}
	em.AddListener(&listener)
	assert.Equal(t, len(em.Listeners[listener.EventName]), 1)
	assert.Equal(t, em.Listeners[listener.EventName][0], &listener)
}

func TestRemovesListenerForEvent(t *testing.T) {
	em := InitEmitter()
	listener := Listener{
		EventName:  mockEventName,
		ListenerFn: func(event *Event) {},
	}
	em.Listeners[listener.EventName] = append(em.Listeners[listener.EventName], &listener)
	assert.Equal(t, len(em.Listeners[listener.EventName]), 1)
	em.RemoveListener(&listener)
	assert.Equal(t, len(em.Listeners[listener.EventName]), 0)
}

func TestMapKeyIsPurgedForEvent(t *testing.T) {
	em := InitEmitter()
	listener1 := Listener{EventName: mockEventName, ListenerFn: func(e *Event) {}}
	listener2 := Listener{EventName: mockEventName, ListenerFn: func(e *Event) {}}
	em.Listeners[listener1.EventName] = append(em.Listeners[listener1.EventName], &listener1)
	em.Listeners[listener2.EventName] = append(em.Listeners[listener2.EventName], &listener2)
	em.RemoveListener(&listener1)
	em.RemoveListener(&listener2)
	assert.NotContains(t, em.Listeners, mockEventName)
}

func TestListsAllRegisteredEvents(t *testing.T) {
	em := InitEmitter()
	listener := Listener{
		EventName:  mockEventName,
		ListenerFn: func(event *Event) {},
	}
	expectedEventNames := []string{mockEventName}
	em.AddListener(&listener)
	assert.EqualValues(t, em.EventNames(), expectedEventNames)
}

type MockCallable struct {
	mock.Mock
	waitGroup *sync.WaitGroup
}

func (m *MockCallable) handleMockEvent(e *Event) {
	defer m.waitGroup.Done()
	m.Called(e.EventName)
}

func TestEventIsEmitted(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	mockCallable := MockCallable{waitGroup: &wg}
	mockCallable.On("handleMockEvent", mock.Anything).Return()
	mockListener := Listener{EventName: mockEventName, ListenerFn: mockCallable.handleMockEvent}
	mockEvent := Event{EventName: mockEventName}
	em := InitEmitter()
	em.AddListener(&mockListener)
	em.Emit(&mockEvent)
	wg.Wait()
	mockCallable.AssertCalled(t, "handleMockEvent", mock.Anything)
}
