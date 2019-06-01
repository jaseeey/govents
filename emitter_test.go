package events

import (
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
)

type MockCallable struct {
	mock.Mock
	waitGroup *sync.WaitGroup
}

func (m *MockCallable) doSomething(data interface{}) {
	defer m.waitGroup.Done()
	m.Called(data)
}

func TestEmitEvent(t *testing.T) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	mockData := "this is just some data"
	mockCallable := MockCallable{waitGroup: &waitGroup}
	mockCallable.On("doSomething", mockData)
	mockListener := Listener{
		eventId:    "some-event",
		callableFn: mockCallable.doSomething,
	}
	AttachListener(&mockListener)
	EmitEvent("some-event", mockData)
	waitGroup.Wait()
	mockCallable.AssertCalled(t, "doSomething", mockData)
}
