package events

import (
	"gotest.tools/assert"
	"reflect"
	"testing"
)

const mockEventId = "do-something"
const mockMessage = "This is a test data for the event"

func TestEventHasIdAndMessage(t *testing.T) {
	mockEvent := Event{
		id:   mockEventId,
		data: mockMessage,
	}
	assert.Equal(t, mockEvent.id, mockEventId)
	assert.Equal(t, mockEvent.data, mockMessage)
}

func TestAddListener(t *testing.T) {
	mockCallableFn := func(data interface{}) {}
	mockListener := Listener{
		eventId:    mockEventId,
		callableFn: mockCallableFn,
	}
	defer ClearAllListeners(mockEventId)
	AttachListener(&mockListener)
	assert.Equal(t, len(allListeners[mockEventId]), 1)
	assert.Equal(t, allListeners[mockEventId][0].eventId, mockEventId)
	assert.Equal(t, reflect.ValueOf(allListeners[mockEventId][0].callableFn).Pointer(), reflect.ValueOf(mockCallableFn).Pointer())
}

func TestDetachListener(t *testing.T) {
	mockCallableFn := func(data interface{}) {}
	mockListener := Listener{
		eventId:    mockEventId,
		callableFn: mockCallableFn,
	}
	defer ClearAllListeners(mockEventId)
	AttachListener(&mockListener)
	assert.Equal(t, len(allListeners[mockEventId]), 1)
	DetachListener(&mockListener)
	assert.Equal(t, len(allListeners[mockEventId]), 0)
}
