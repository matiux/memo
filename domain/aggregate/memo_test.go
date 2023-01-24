package aggregate

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_it_should_be_created(t *testing.T) {

	memoId := NewUUIDv4()
	body := "Vegetables are good"
	creationDate := time.Now()

	memo := NewMemo(memoId, body, creationDate)
	//fmt.Printf("%v\n", memoId)
	//fmt.Printf("%v\n", memo.id)
	assert.True(t, memo.id.Equals(memoId))
	assert.Equal(t, body, memo.body)
	assert.Equal(t, creationDate, memo.creationDate)
	assert.Len(t, memo.UncommittedEvents, 1)

	//id := "message_id"
	//playhead := 1
	//metadata := Metadata{MetadataValuesT{"foo": "bar"}}
	//eventType := "OrderCreated"
	//payload := []byte(`{"some":"json"}`)
	//
	//domainMessage := RecordMessageNow(id, playhead, metadata, eventType, payload)
	//
	//assert.Equal(t, id, domainMessage.Id)
	//assert.Equal(t, playhead, domainMessage.Playhead)
	//assert.Equal(t, metadata, domainMessage.Metadata)
	//assert.Equal(t, eventType, domainMessage.EventType)
	//assert.Equal(t, payload, domainMessage.Payload)

}
