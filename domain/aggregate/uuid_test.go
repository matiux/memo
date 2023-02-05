package aggregate_test

import (
	"github.com/matiux/memo/domain/aggregate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUID_it_should_be_true_when_two_id_equals(t *testing.T) {

	stringUUID := "f9bee14a-c795-4fc9-8e45-e0fa1759f347"

	id1 := aggregate.NewUUIDv4From(stringUUID)
	id2 := aggregate.NewUUIDv4From(stringUUID)

	assert.True(t, id1.Equals(id2))
}
