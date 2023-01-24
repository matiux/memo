package tmp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_it_has_getters(t *testing.T) {

	id := "message_id"
	playhead := 1
	metadata := Metadata{MetadataValuesT{"foo": "bar"}}
	eventType := "OrderCreated"
	payload := []byte(`{"some":"json"}`)

	domainMessage := RecordMessageNow(id, playhead, metadata, eventType, payload)

	assert.Equal(t, id, domainMessage.Id)
	assert.Equal(t, playhead, domainMessage.Playhead)
	assert.Equal(t, metadata, domainMessage.Metadata)
	assert.Equal(t, eventType, domainMessage.EventType)
	assert.Equal(t, payload, domainMessage.Payload)

}

func Test_it_returns_a_new_instance_with_more_metadata_on_and_metadata(t *testing.T) {

	domainMessage := RecordMessageNow(
		"message_id",
		1,
		NewMetadata(MetadataValuesT{}),
		"OrderCreated",
		[]byte(`{"some":"json"}`),
	)
	newDomainMessage := domainMessage.andMetadata(NewMetadataKV("key", "value"))

	assert.NotEqual(t, domainMessage, newDomainMessage)

	assert.Len(t, domainMessage.Metadata.values, 0)
	assert.Len(t, newDomainMessage.Metadata.values, 1)
}

//func Test_it_keeps_all_data_the_same_expect_metadata_on_and_metadata(t *testing.T) {
//
//	domainMessage := RecordMessageNow("message_id", 42, NewMetadata(MetadataValuesT{}), "payload")
//	newDomainMessage := domainMessage.andMetadata(NewMetadataKV("key", "value"))
//
//	assert.Equal(t, domainMessage.id, newDomainMessage.id)
//	assert.Equal(t, domainMessage.playhead, newDomainMessage.playhead)
//	assert.Equal(t, domainMessage.payload, newDomainMessage.payload)
//	assert.Equal(t, domainMessage.recordedOn, newDomainMessage.recordedOn)
//
//	assert.NotEqual(t, domainMessage.metadata, newDomainMessage.metadata)
//
//}
//
//func Test_it_merges_the_metadata_instances_on_and_metadata(t *testing.T) {
//
//	domainMessage := RecordMessageNow("message_id", 42, NewMetadataKV("key", "value"), "payload").
//		andMetadata(NewMetadataKV("foo", "bar"))
//
//	expected := Metadata{MetadataValuesT{"key": "value", "foo": "bar"}}
//
//	assert.Equal(t, expected, domainMessage.metadata)
//}
