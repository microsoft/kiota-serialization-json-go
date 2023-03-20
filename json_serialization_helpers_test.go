package jsonserialization

import (
	testing "testing"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-serialization-json-go/internal"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		name       string
		useDefault bool
	}{
		{name: "uses default ParseNode", useDefault: true},
		{name: "uses fallback", useDefault: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.useDefault {
				absser.DefaultParseNodeFactoryInstance.ContentTypeAssociatedFactories["application/json"] = NewJsonParseNodeFactory()
				defer delete(absser.DefaultParseNodeFactoryInstance.ContentTypeAssociatedFactories, "application/json")
			}

			source := "{\"displayName\":\"McGill\",\"officeLocation\":\"Montreal\", \"id\": \"opaque\"}"
			sourceArray := []byte(source)

			result := internal.NewIntersectionTypeMock()
			err := Unmarshal(sourceArray, &result, internal.CreateIntersectionTypeMockFromDiscriminator)

			if err != nil {
				t.Error(err)
			}

			assert.NotNil(t, result)
			assert.NotNil(t, result.GetComposedType1())
			assert.NotNil(t, result.GetComposedType2())
			assert.Nil(t, result.GetStringValue())
			assert.Nil(t, result.GetComposedType3())
			assert.Equal(t, "McGill", *result.GetComposedType2().GetDisplayName())
			assert.Equal(t, "opaque", *result.GetComposedType1().GetId())
		})
	}
}

func TestUnmarshalFromNull(t *testing.T) {
	sourceArray := []byte("null")

	result := internal.NewIntersectionTypeMock()
	err := Unmarshal(sourceArray, &result, internal.CreateIntersectionTypeMockFromDiscriminator)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestUnmarshalWithError(t *testing.T) {
	sourceArray := []byte("}")

	result := internal.NewIntersectionTypeMock()
	err := Unmarshal(sourceArray, &result, internal.CreateIntersectionTypeMockFromDiscriminator)
	assert.Error(t, err)
}

func TestMarshal(t *testing.T) {
	testCases := []struct {
		name       string
		useDefault bool
	}{
		{name: "uses default ParseNodeFactory", useDefault: true},
		{name: "uses fallback", useDefault: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.useDefault {
				absser.DefaultParseNodeFactoryInstance.ContentTypeAssociatedFactories["application/json"] = NewJsonParseNodeFactory()
				defer delete(absser.DefaultParseNodeFactoryInstance.ContentTypeAssociatedFactories, "application/json")
				absser.DefaultSerializationWriterFactoryInstance.ContentTypeAssociatedFactories["application/json"] = NewJsonSerializationWriterFactory()
				defer delete(absser.DefaultSerializationWriterFactoryInstance.ContentTypeAssociatedFactories, "application/json")
			}

			source := "{\"displayName\":\"McGill\",\"officeLocation\":\"Montreal\", \"id\": \"opaque\"}"
			sourceArray := []byte(source)

			result := internal.NewIntersectionTypeMock()
			err := Unmarshal(sourceArray, &result, internal.CreateIntersectionTypeMockFromDiscriminator)

			if err != nil {
				t.Error(err)
			}

			b, err := Marshal(result)
			assert.NoError(t, err)
			assert.JSONEq(t, source, string(b))
		})
	}
}

func TestMarshalToNull(t *testing.T) {
	b, err := Marshal(nil)
	assert.NoError(t, err)
	assert.JSONEq(t, "null", string(b))

	var (
		entity    *internal.TestEntity
		zeroValue absser.Parsable = entity
	)

	b, err = Marshal(zeroValue)
	assert.NoError(t, err)
	assert.JSONEq(t, "null", string(b))
}
