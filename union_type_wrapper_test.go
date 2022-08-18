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

func TestItParsesUnionTypeStringValue(t *testing.T) {
	source := "\"officeLocation\""
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
	assert.Nil(t, cast.GetComposedType3())
	assert.NotNil(t, cast.GetStringValue())
	assert.Equal(t, "officeLocation", *cast.GetStringValue())
}

func TestItSerializesUnionTypeStringValue(t *testing.T) {
	value := "officeLocation"
	source := internal.NewUnionTypeMockable()
	source.SetStringValue(&value)
	writer := NewJsonSerializationWriter()
	err := source.Serialize(writer)
	if err != nil {
		t.Error(err)
	}
	result, err := writer.GetSerializedContent()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "\"officeLocation\"", string(result))
	defer writer.Close()
}

func TestItSerializesUnionTypeComplexProperty1(t *testing.T) {
	source := internal.NewUnionTypeMockable()
	prop1Value := internal.NewTestEntity()
	idValue := "opaque"
	prop1Value.SetId(&idValue)
	officeLocationValue := "Montreal"
	prop1Value.SetOfficeLocation(&officeLocationValue)
	prop2Value := internal.NewSecondTestEntity()
	idIntValue := int64(10)
	prop2Value.SetId(&idIntValue)
	displayNameValue := "McGill"
	prop2Value.SetDisplayName(&displayNameValue)
	source.SetComposedType1(prop1Value)
	source.SetComposedType2(prop2Value)
	writer := NewJsonSerializationWriter()
	err := source.Serialize(writer)
	if err != nil {
		t.Error(err)
	}
	result, err := writer.GetSerializedContent()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "{\"id\":\"opaque\",\"officeLocation\":\"Montreal\"}", string(result))
	defer writer.Close()
}

func TestItSerializesUnionTypeComplexProperty2(t *testing.T) {
	source := internal.NewUnionTypeMockable()
	prop2Value := internal.NewSecondTestEntity()
	idIntValue := int64(10)
	prop2Value.SetId(&idIntValue)
	displayNameValue := "McGill"
	prop2Value.SetDisplayName(&displayNameValue)
	source.SetComposedType2(prop2Value)
	writer := NewJsonSerializationWriter()
	err := source.Serialize(writer)
	if err != nil {
		t.Error(err)
	}
	result, err := writer.GetSerializedContent()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "{\"id\":10,\"displayName\":\"McGill\"}", string(result))
	defer writer.Close()
}

func TestItSerializesUnionTypeComplexProperty3(t *testing.T) {
	source := internal.NewUnionTypeMockable()
	prop3Value1 := internal.NewTestEntity()
	idIntValue := "10"
	prop3Value1.SetId(&idIntValue)
	officeLocationValue1 := "Montreal"
	prop3Value1.SetOfficeLocation(&officeLocationValue1)
	prop3Value2 := internal.NewTestEntity()
	idIntValue2 := "11"
	prop3Value2.SetId(&idIntValue2)
	officeLocationValue2 := "Ottawa"
	prop3Value2.SetOfficeLocation(&officeLocationValue2)
	source.SetComposedType3([]internal.TestEntityable{prop3Value1, prop3Value2})
	writer := NewJsonSerializationWriter()
	err := source.Serialize(writer)
	if err != nil {
		t.Error(err)
	}
	result, err := writer.GetSerializedContent()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[{\"id\":\"10\",\"officeLocation\":\"Montreal\"},{\"id\":\"11\",\"officeLocation\":\"Ottawa\"}]", string(result))
	defer writer.Close()
}
