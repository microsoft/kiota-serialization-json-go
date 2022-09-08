package jsonserialization

import (
	"testing"

	"github.com/microsoft/kiota-serialization-json-go/internal"
	"github.com/stretchr/testify/assert"
)

func TestItParsesIntersectionTypeComplexProperty1(t *testing.T) {
	source := "{\"displayName\":\"McGill\",\"officeLocation\":\"Montreal\", \"id\": \"opaque\"}"
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)

	if err != nil {
		t.Error(err)
	}
	result, err := parseNode.GetObjectValue(internal.CreateIntersectionTypeMockFromDiscriminator)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
	cast, ok := result.(internal.IntersectionTypeMockable)
	assert.True(t, ok)
	assert.NotNil(t, cast.GetComposedType1())
	assert.NotNil(t, cast.GetComposedType2())
	assert.Nil(t, cast.GetStringValue())
	assert.Nil(t, cast.GetComposedType3())
	assert.Equal(t, "McGill", *cast.GetComposedType2().GetDisplayName())
	assert.Equal(t, "opaque", *cast.GetComposedType1().GetId())
}

func TestItParsesIntersectionTypeComplexProperty2(t *testing.T) {
	source := "{\"displayName\":\"McGill\",\"officeLocation\":\"Montreal\", \"id\": 10}"
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)

	if err != nil {
		t.Error(err)
	}
	result, err := parseNode.GetObjectValue(internal.CreateIntersectionTypeMockFromDiscriminator)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
	cast, ok := result.(internal.IntersectionTypeMockable)
	assert.True(t, ok)
	assert.NotNil(t, cast.GetComposedType1())
	assert.NotNil(t, cast.GetComposedType2())
	assert.Nil(t, cast.GetStringValue())
	assert.Nil(t, cast.GetComposedType3())
	assert.Nil(t, cast.GetComposedType1().GetId())
	assert.Nil(t, cast.GetComposedType2().GetId()) // it's expected to be null since we have conflicting properties here and the parser will only try one to avoid having to brute its way through
	assert.Equal(t, "McGill", *cast.GetComposedType2().GetDisplayName())
}

func TestItParsesIntersectionTypeComplexProperty3(t *testing.T) {
	source := "[{\"@odata.type\":\"#microsoft.graph.TestEntity\",\"officeLocation\":\"Ottawa\", \"id\": \"11\"}, {\"@odata.type\":\"#microsoft.graph.TestEntity\",\"officeLocation\":\"Montreal\", \"id\": \"10\"}]"
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)

	if err != nil {
		t.Error(err)
	}
	result, err := parseNode.GetObjectValue(internal.CreateIntersectionTypeMockFromDiscriminator)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
	cast, ok := result.(internal.IntersectionTypeMockable)
	assert.True(t, ok)
	assert.Nil(t, cast.GetComposedType1())
	assert.Nil(t, cast.GetComposedType2())
	assert.Nil(t, cast.GetStringValue())
	assert.NotNil(t, cast.GetComposedType3())
	assert.Equal(t, 2, len(cast.GetComposedType3()))
	assert.Equal(t, "Ottawa", *cast.GetComposedType3()[0].GetOfficeLocation())
}

func TestItParsesIntersectionTypeStringValue(t *testing.T) {
	source := "\"officeLocation\""
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)

	if err != nil {
		t.Error(err)
	}
	result, err := parseNode.GetObjectValue(internal.CreateIntersectionTypeMockFromDiscriminator)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
	cast, ok := result.(internal.IntersectionTypeMockable)
	assert.True(t, ok)
	assert.Nil(t, cast.GetComposedType1())
	assert.Nil(t, cast.GetComposedType2())
	assert.NotNil(t, cast.GetStringValue())
	assert.Nil(t, cast.GetComposedType3())
	assert.Equal(t, "officeLocation", *cast.GetStringValue())
}

func TestItSerializesIntersectionTypeComplexProperty1(t *testing.T) {
	source := internal.NewIntersectionTypeMock()
	prop1Value := internal.NewTestEntity()
	idValue := "opaque"
	prop1Value.SetId(&idValue)
	officeLocationValue := "Montreal"
	prop1Value.SetOfficeLocation(&officeLocationValue)
	prop2Value := internal.NewSecondTestEntity()
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
	assert.Equal(t, "{\"id\":\"opaque\",\"officeLocation\":\"Montreal\",\"displayName\":\"McGill\"}", string(result))
	defer writer.Close()
}

func TestItSerializesIntersectionTypeComplexProperty2(t *testing.T) {
	source := internal.NewIntersectionTypeMock()
	prop2Value := internal.NewSecondTestEntity()
	displayNameValue := "McGill"
	prop2Value.SetDisplayName(&displayNameValue)
	idValue := int64(10)
	prop2Value.SetId(&idValue)
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

func TestItSerializesIntersectionTypeComplexProperty(t *testing.T) {
	source := internal.NewIntersectionTypeMock()
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
