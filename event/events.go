package event

var bus = NewEventGroup()

func RegisterSubscribe(topic string, listener any) error {
	return bus.RegisterSubscribe(topic, listener)
}

func RegisterAsyncSubscribe(topic string, listener any) error {
	return bus.RegisterAsyncSubscribe(topic, listener)
}

func UnregisterSubscribe(topic string, listener any) error {
	return bus.UnregisterSubscribe(topic, listener)
}

func Publish(topic string, args ...any) {
	bus.Publish(topic, args...)
}
