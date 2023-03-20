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

func TestInvalidContentShouldFail(t *testing.T) {
	source := "3 [ }"
	sourceArray := []byte(source)

	parseNode, err := NewJsonParseNode(sourceArray)
	assert.Error(t, err)
	assert.Nil(t, parseNode)
}
