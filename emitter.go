package events

func EmitEvent(eventName string, data interface{}) {
	listeners := getAllListeners(eventName)
	for _, listener := range listeners {
		go listener.callableFn(data)
	}
}

