package domain_test

import (
	"fmt"
	"github.com/matiux/memo/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var eventBus domain.EventBus

func setupTestEventBus() {
	eventBus = domain.NewSimpleEventBus()
}

func createTestDomainMessage(body string) domain.Message {

	event := &eventOccurred{domain.NewUUIDv4(), body, domain.BasicEvent{}}

	return domain.Message{
		Playhead:    domain.Playhead(1),
		EventType:   event.Kind(),
		Payload:     event,
		AggregateId: domain.NewUUIDv4(),
		RecordedOn:  time.Now(),
	}
}

func TestEventBus_it_subscribes_an_event_listener(t *testing.T) {

	setupTestEventBus()

	domainMessage := createTestDomainMessage("The event Body")
	eventListener := &eventListenerMock{}
	eventListener.On("Handle", domainMessage).Once().Return(nil)

	eventBus.Subscribe(eventListener)
	err := eventBus.Publish(domain.EventStream{domainMessage})

	assert.Nil(t, err)
	eventListener.AssertExpectations(t)
}

func TestEventBus_it_publishes_events_to_subscribed_event_listeners(t *testing.T) {

	setupTestEventBus()

	domainMessage1 := createTestDomainMessage("The event Body 1")
	domainMessage2 := createTestDomainMessage("The event Body 2")

	domainEventStream := domain.EventStream{domainMessage1, domainMessage2}

	eventListener1 := &eventListenerMock{}
	eventListener1.On("Handle", domainMessage1).Once().Return(nil)
	eventListener1.On("Handle", domainMessage2).Once().Return(nil)

	eventListener2 := &eventListenerMock{}
	eventListener2.On("Handle", domainMessage1).Once().Return(nil)
	eventListener2.On("Handle", domainMessage2).Once().Return(nil)

	eventBus.Subscribe(eventListener1)
	eventBus.Subscribe(eventListener2)
	err := eventBus.Publish(domainEventStream)

	assert.Nil(t, err)
	eventListener1.AssertExpectations(t)
	eventListener2.AssertExpectations(t)
}

func TestEventBus_it_does_not_dispatch_new_events_before_all_listeners_have_run(t *testing.T) {

	setupTestEventBus()

	domainMessage1 := createTestDomainMessage("The event Body 1")
	domainMessage2 := createTestDomainMessage("The event Body 2")

	domainEventStream := domain.EventStream{domainMessage1}

	eventListener1 := simpleEventBusTestListener{
		eventBus,
		domain.EventStream{domainMessage2},
		false,
	}

	eventListener2 := &eventListenerMock{}
	eventListener2.On("Handle", domainMessage1).Once().Return(nil)
	eventListener2.On("Handle", domainMessage2).Once().Return(nil)

	eventBus.Subscribe(&eventListener1)
	eventBus.Subscribe(eventListener2)
	err := eventBus.Publish(domainEventStream)

	assert.Nil(t, err)
	eventListener2.AssertExpectations(t)
}

func TestEventBus_it_should_still_publish_events_after_exception(t *testing.T) {

	setupTestEventBus()

	domainMessage1 := createTestDomainMessage("The event Body 1")
	domainMessage2 := createTestDomainMessage("The event Body 2")

	domainEventStream1 := domain.EventStream{domainMessage1}
	domainEventStream2 := domain.EventStream{domainMessage2}

	eventListener := &eventListenerMock{}
	eventListener.On("Handle", domainMessage1).Once().Return(fmt.Errorf("an error"))
	eventListener.On("Handle", domainMessage2).Once().Return(nil)

	eventBus.Subscribe(eventListener)

	err := eventBus.Publish(domainEventStream1)

	assert.NotNil(t, err)
	assert.Equal(t, "Error in Payload Listener `eventListenerMock` with Message `eventOccurred`. Original error: an error", err.Error())

	_ = eventBus.Publish(domainEventStream2)

	eventListener.AssertExpectations(t)
}

type eventOccurred struct {
	id   domain.UUIDv4
	body string
	domain.BasicEvent
}

func (e eventOccurred) Kind() string {
	return "eventOccurred"
}

func (e eventOccurred) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (e eventOccurred) UnmarshalJSON(b []byte) error {

	return nil
}

type simpleEventBusTestListener struct {
	domain.EventBus
	publishableStream domain.EventStream
	handled           bool
}

func (eb *simpleEventBusTestListener) Handle(message domain.Message) error {

	if !eb.handled {
		eb.EventBus.Publish(eb.publishableStream)
		eb.handled = true
	}

	return nil
}

type eventListenerMock struct {
	mock.Mock
}

func (m *eventListenerMock) Handle(message domain.Message) error {
	args := m.Called(message)

	return args.Error(0)
}
