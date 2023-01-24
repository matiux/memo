package tmp

type MetadataValuesT map[string]interface{}

// Metadata adding extra information to the DomainMessage.
type Metadata struct {
	values MetadataValuesT
}

// Merges the values of this and the other instance.
func (m Metadata) merge(otherMetadata Metadata) Metadata {

	values := MetadataValuesT{}

	for k, v := range m.values {
		values[k] = v
	}

	for k, v := range otherMetadata.values {
		values[k] = v
	}

	return Metadata{values}
}

// Get a specific metadata value based on key.
func (m Metadata) get(key string) interface{} {
	return m.values[key]
}

// NewMetadata is the constructor
func NewMetadata(metadata MetadataValuesT) Metadata {
	return Metadata{metadata}
}

// NewMetadataKV is a helper method to construct an instance containing the key and value.
func NewMetadataKV(key string, value interface{}) Metadata {
	return Metadata{MetadataValuesT{key: value}}
}
