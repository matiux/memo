package tmp

import (
	"time"
)

// Message represents an important change in the domain.
type DomainMessage struct {
	Playhead   int
	Metadata   Metadata
	EventType  string
	Payload    []byte
	Id         string
	RecordedOn time.Time
}

func (dm DomainMessage) andMetadata(metadata Metadata) DomainMessage {
	newMetadata := dm.Metadata.merge(metadata)

	return DomainMessage{
		dm.Playhead,
		newMetadata,
		dm.EventType,
		dm.Payload,
		dm.Id,
		dm.RecordedOn,
	}
}

//func (dm *Message) recordNow(id string, playhead int, metadata Metadata, payload interface{}) *Message {
//	return &Message{playhead, metadata, payload, id, time.Now()}
//}

func RecordMessageNow(id string, playhead int, metadata Metadata, eventType string, payload []byte) DomainMessage {
	return DomainMessage{
		playhead,
		metadata,
		eventType,
		payload,
		id,
		time.Now(),
	}
}
