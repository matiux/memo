package eventstore

import (
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var eventBus EventBus

type EventOccurred struct {
	id   aggregate.UUIDv4
	body string
	aggregate.BasicEvent
}

func (e EventOccurred) Kind() string {
	return "EventOccurred"
}

func setupTestEventBus() {
	eventBus = NewSimpleEventBus()
}

type EventListenerMock struct {
	mock.Mock
}

func (m *EventListenerMock) handle(message aggregate.DomainMessage) error {
	args := m.Called(message)

	return args.Error(0)
}

func createTestDomainMessage(body string) aggregate.DomainMessage {

	event := EventOccurred{aggregate.NewUUIDv4(), body, aggregate.BasicEvent{}}

	return aggregate.DomainMessage{
		Playhead:    aggregate.Playhead(1),
		EventType:   event.Kind(),
		Event:       event,
		AggregateId: aggregateId,
		RecordedOn:  time.Now(),
	}
}

func TestEventBus_it_subscribes_an_event_listener(t *testing.T) {

	setupTestEventBus()

	domainMessage := createTestDomainMessage("The event body")
	eventListener := &EventListenerMock{}
	eventListener.On("handle", domainMessage).Once().Return(nil)

	eventBus.subscribe(eventListener)
	eventBus.publish(aggregate.DomainEventStream{domainMessage})

	eventListener.AssertExpectations(t)
}

func TestEventBus_it_publishes_events_to_subscribed_event_listeners(t *testing.T) {

	setupTestEventBus()

	domainMessage1 := createTestDomainMessage("The event body 1")
	domainMessage2 := createTestDomainMessage("The event body 2")

	domainEventStream := aggregate.DomainEventStream{domainMessage1, domainMessage2}

	eventListener1 := &EventListenerMock{}
	eventListener1.On("handle", domainMessage1).Once().Return(nil)
	eventListener1.On("handle", domainMessage2).Once().Return(nil)

	eventListener2 := &EventListenerMock{}
	eventListener2.On("handle", domainMessage1).Once().Return(nil)
	eventListener2.On("handle", domainMessage2).Once().Return(nil)

	eventBus.subscribe(eventListener1)
	eventBus.subscribe(eventListener2)
	eventBus.publish(domainEventStream)

	eventListener1.AssertExpectations(t)
	eventListener2.AssertExpectations(t)
}

type SimpleEventBusTestListener struct {
	EventBus
	publishableStream aggregate.DomainEventStream
	handled           bool
}

func (eb *SimpleEventBusTestListener) handle(message aggregate.DomainMessage) error {

	if !eb.handled {
		eb.EventBus.publish(eb.publishableStream)
		eb.handled = true
	}

	return nil
}

func TestEventBus_it_does_not_dispatch_new_events_before_all_listeners_have_run(t *testing.T) {

	setupTestEventBus()

	domainMessage1 := createTestDomainMessage("The event body 1")
	domainMessage2 := createTestDomainMessage("The event body 2")

	domainEventStream := aggregate.DomainEventStream{domainMessage1}

	eventListener1 := SimpleEventBusTestListener{
		eventBus,
		aggregate.DomainEventStream{domainMessage2},
		false,
	}

	eventListener2 := &EventListenerMock{}
	eventListener2.On("handle", domainMessage1).Once().Return(nil)
	eventListener2.On("handle", domainMessage2).Once().Return(nil)

	eventBus.subscribe(&eventListener1)
	eventBus.subscribe(eventListener2)
	eventBus.publish(domainEventStream)

	eventListener2.AssertExpectations(t)
}

func TestEventBus_it_should_still_publish_events_after_exception(t *testing.T) {

	setupTestEventBus()

	domainMessage1 := createTestDomainMessage("The event body 1")
	domainMessage2 := createTestDomainMessage("The event body 2")

	domainEventStream1 := aggregate.DomainEventStream{domainMessage1}
	domainEventStream2 := aggregate.DomainEventStream{domainMessage2}

	eventListener := &EventListenerMock{}
	eventListener.On("handle", domainMessage1).Once().Return(fmt.Errorf("an error"))
	eventListener.On("handle", domainMessage2).Once().Return(nil)

	eventBus.subscribe(eventListener)

	err := eventBus.publish(domainEventStream1)

	assert.NotNil(t, err)
	assert.Equal(t, "Error in Event Listener `EventListenerMock` with Message `EventOccurred`. Original error: an error", err.Error())

	_ = eventBus.publish(domainEventStream2)

	eventListener.AssertExpectations(t)
}
