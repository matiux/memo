package tmp

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_it_contains_values_from_both_instances_after_merge(t *testing.T) {

	metadata1 := NewMetadata(MetadataValuesT{"foo": 42})
	metadata2 := NewMetadataKV("bar", "value")

	expectedMetadata := NewMetadata(MetadataValuesT{"foo": 42, "bar": "value"})

	assert.True(t, reflect.DeepEqual(expectedMetadata, metadata1.merge(metadata2)))
}

func Test_it_overrides_values_with_data_from_other_instance_on_merge(t *testing.T) {

	metadata1 := NewMetadataKV("foo", "value")
	metadata2 := NewMetadata(MetadataValuesT{"foo": 42})

	expectedMetadata := NewMetadata(MetadataValuesT{"foo": 42})

	assert.True(t, reflect.DeepEqual(expectedMetadata, metadata1.merge(metadata2)), "True is true!")
}

func Test_it_constructs_an_instance_containing_the_key_and_value(t *testing.T) {

	metadata := NewMetadataKV("foo", 42)
	expectedMetadata := NewMetadata(MetadataValuesT{"foo": 42})

	assert.Equal(t, expectedMetadata, metadata)
}

func Test_it_returns_all_values(t *testing.T) {

	metadata := NewMetadata(MetadataValuesT{"foo": 42, "bar": "value"})

	expected := MetadataValuesT{"foo": 42, "bar": "value"}

	assert.Equal(t, expected, metadata.values)
}

func Test_it_returns_nil_when_get_contains_unset_key(t *testing.T) {

	metadata := NewMetadata(MetadataValuesT{"foo": 42})

	assert.Nil(t, metadata.get("bar"))
}

func Test_it_returns_the_value_of_a_key_with_get(t *testing.T) {

	metadata := NewMetadata(MetadataValuesT{"foo": 42})

	assert.Equal(t, 42, metadata.get("foo"))
}
