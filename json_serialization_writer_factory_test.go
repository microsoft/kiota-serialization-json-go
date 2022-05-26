package jsonserialization

import (
	"testing"

	assert "github.com/stretchr/testify/assert"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

func TestJsonSerializationFactoryWriterHonoursInterface(t *testing.T) {
	instance := NewJsonSerializationWriterFactory()
	assert.Implements(t, (*absser.SerializationWriterFactory)(nil), instance)
}
