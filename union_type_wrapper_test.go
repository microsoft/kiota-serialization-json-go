package jsonserialization

import (
	"testing"

	"github.com/microsoft/kiota-serialization-json-go/internal"
	"github.com/stretchr/testify/assert"
)

func TestItParsesUnionTypeComplexProperty1(t *testing.T) {
	source := "{\"@odata.type\":\"#microsoft.graph.testEntity\",\"officeLocation\":\"Montreal\", \"id\": \"opaque\"}"
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)

	if err != nil {
		t.Error(err)
	}
	result, err := parseNode.GetObjectValue(internal.CreateUnionTypeMockableFromDiscriminator)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
	cast, ok := result.(internal.UnionTypeMockable)
	assert.True(t, ok)
	assert.NotNil(t, cast.GetComposedType1())
	assert.Nil(t, cast.GetComposedType2())
	assert.Nil(t, cast.GetStringValue())
	assert.Nil(t, cast.GetComposedType3())
	assert.Equal(t, "Montreal", *cast.GetComposedType1().GetOfficeLocation())
	assert.Equal(t, "opaque", *cast.GetComposedType1().GetId())
}

func TestItParsesUnionTypeComplexProperty2(t *testing.T) {
	source := "{\"@odata.type\":\"#microsoft.graph.secondTestEntity\",\"officeLocation\":\"Montreal\", \"id\": 10}"
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)

	if err != nil {
		t.Error(err)
	}
	result, err := parseNode.GetObjectValue(internal.CreateUnionTypeMockableFromDiscriminator)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
	cast, ok := result.(internal.UnionTypeMockable)
	assert.True(t, ok)
	assert.Nil(t, cast.GetComposedType1())
	assert.NotNil(t, cast.GetComposedType2())
	assert.Nil(t, cast.GetStringValue())
	assert.Nil(t, cast.GetComposedType3())
	assert.Equal(t, int64(10), *cast.GetComposedType2().GetId())
}

func TestItParsesUnionTypeComplexProperty3(t *testing.T) {
	source := "[{\"@odata.type\":\"#microsoft.graph.TestEntity\",\"officeLocation\":\"Ottawa\", \"id\": \"11\"}, {\"@odata.type\":\"#microsoft.graph.TestEntity\",\"officeLocation\":\"Montreal\", \"id\": \"10\"}]"
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)

	if err != nil {
		t.Error(err)
	}
	result, err := parseNode.GetObjectValue(internal.CreateUnionTypeMockableFromDiscriminator)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
	cast, ok := result.(internal.UnionTypeMockable)
	assert.True(t, ok)
	assert.Nil(t, cast.GetComposedType1())
	assert.Nil(t, cast.GetComposedType2())
	assert.Nil(t, cast.GetStringValue())
	assert.NotNil(t, cast.GetComposedType3())
	assert.Equal(t, 2, len(cast.GetComposedType3()))
	assert.Equal(t, "11", *cast.GetComposedType3()[0].GetId())
}
