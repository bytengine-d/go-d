package event

import "github.com/asaskevich/EventBus"

type EventGroup struct {
	eventBus EventBus.Bus
}

func (ep *EventGroup) RegisterSubscribe(topic string, listener any) error {
	return ep.eventBus.Subscribe(topic, listener)
}

func (ep *EventGroup) RegisterAsyncSubscribe(topic string, listener any) error {
	return ep.eventBus.SubscribeAsync(topic, listener, false)
}

func (ep *EventGroup) UnregisterSubscribe(topic string, listener any) error {
	return ep.eventBus.Unsubscribe(topic, listener)
}

func (ep *EventGroup) Publish(topic string, args ...any) {
	ep.eventBus.Publish(topic, args...)
}

func NewEventGroup() *EventGroup {
	return &EventGroup{
		eventBus: EventBus.New(),
	}
}
