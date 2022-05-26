package jsonserialization

import (
	"testing"

	assert "github.com/stretchr/testify/assert"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

func TestJsonParseNodeFactoryHonoursInterface(t *testing.T) {
	instance := NewJsonParseNodeFactory()
	assert.Implements(t, (*absser.ParseNodeFactory)(nil), instance)
}
