package variant_test

import (
	"github.com/matiux/memo/domain/variant"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUIDVariant_it_should_be_true_when_two_id_equals(t *testing.T) {

	id1 := variant.NewUuidV4("2022")
	id2 := variant.NewUuidV4("2022")

	assert.True(t, id1.Equals(id2))
	assert.True(t, variant.UuidV4{}.Create("2022").Equals(id2))

	id3 := variant.NewIntId(2022)
	id4 := variant.NewIntId(2022)

	assert.True(t, id3.Equals(id4))
	assert.True(t, variant.IntId{}.Create(2022).Equals(id4))
}
