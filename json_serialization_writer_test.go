package jsonserialization

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/microsoft/kiota-serialization-json-go/internal"

	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

func referenceTime() (value time.Time) {
	value, _ = time.Parse(time.Layout, time.Layout)
	return
}

func TestItDoesntWriteAnythingForNilAdditionalData(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	serializer.WriteAdditionalData(nil)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestItDoesntWriteAnythingForEmptyAdditionalData(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	serializer.WriteAdditionalData(make(map[string]interface{}))
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestItDoesntTrimCommasOnEmptyAdditionalData(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "value"
	serializer.WriteStringValue("key", &value)
	serializer.WriteAdditionalData(make(map[string]interface{}))
	value2 := "value2"
	serializer.WriteStringValue("key2", &value2)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, "\"key\":\"value\",\"key2\":\"value2\"", string(result[:]))
}

func TestWriteTimeValue(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := referenceTime()
	serializer.WriteTimeValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":%q", value.Format(time.RFC3339)), string(result[:]))
}

func TestWriteISODurationValue(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := absser.NewDuration(1, 0, 2, 3, 4, 5, 6)
	serializer.WriteISODurationValue("key", value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":%q", value), string(result[:]))
}

func TestWriteTimeOnlyValue(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := absser.NewTimeOnly(referenceTime())
	serializer.WriteTimeOnlyValue("key", value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":%q", value), string(result[:]))
}

func TestWriteDateOnlyValue(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := absser.NewDateOnly(referenceTime())
	serializer.WriteDateOnlyValue("key", value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":%q", value), string(result[:]))
}

func TestWriteBoolValue(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := true
	serializer.WriteBoolValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":%t", value), string(result[:]))
}

func TestWriteInt8Value(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := int8(125)
	serializer.WriteInt8Value("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":%d", value), string(result[:]))
}

func TestWriteByteValue(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	var value byte = 97
	serializer.WriteByteValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":%d", value), string(result[:]))
}

// ByteArray values are encoded to Base64 when stored
func TestWriteByteArrayValue(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := []byte("SerialWriter")
	serializer.WriteByteArrayValue("key", value)
	expected := "U2VyaWFsV3JpdGVy"
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":\"%s\"", expected), string(result[:]))
}

func TestDoubleEscapeFailure(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "W/\"CQAAABYAAAAs+XSiyjZdS4Rhtwk0v1pGAAC5bsJ2\""
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"key\":%q", value), string(result[:]))
}

func TestReset(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "W/\"CQAAABYAAAAs+XSiyjZdS4Rhtwk0v1pGAAC5bsJ2\""
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.True(t, len(result) > 0)
	serializer.Reset()
	assert.True(t, len(result) > 0)
	empty, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.True(t, len(empty) == 0)
	dateOnly := absser.NewDateOnly(referenceTime())
	serializer.WriteDateOnlyValue("today", dateOnly)
	notEmpty, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.True(t, len(notEmpty) > 0)
}

func TestClose(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	serializer.Close()
	assert.Panics(t, func() {
		serializer.GetSerializedContent()
	})
	assert.Panics(t, serializer.writer.Reset)
	assert.NotPanics(t, func() {
		serializer.Close()
	})
}

func TestJsonSerializationWriterHonoursInterface(t *testing.T) {
	instance := NewJsonSerializationWriter()
	assert.Implements(t, (*absser.SerializationWriter)(nil), instance)
}

func TestWriteMultipleTypes(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "value"
	serializer.WriteStringValue("key", &value)
	pointer := "pointer"
	adlData := map[string]interface{}{
		"add1": "string",
		"add2": &pointer,
		"add3": []string{"foo", "bar"},
	}
	serializer.WriteAdditionalData(adlData)
	value2 := "value2"
	serializer.WriteStringValue("key2", &value2)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Contains(t, string(result[:]), "\"key\":\"value\",")
	assert.Contains(t, string(result[:]), "\"add1\":\"string\",")
	assert.Contains(t, string(result[:]), "\"add2\":\"pointer\",")
	assert.Contains(t, string(result[:]), "\"add3\":[\"foo\",\"bar\"],")
	assert.Contains(t, string(result[:]), "\"key2\":\"value2\"")
	assert.Equal(t, len("\"key\":\"value\",\"add1\":\"string\",\"add2\":\"pointer\",\"add3\":[\"foo\",\"bar\"],\"key2\":\"value2\""), len(string(result[:])))
}

func TestWriteInvalidAdditionalData(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "value"
	serializer.WriteStringValue("key", &value)

	adlData := map[string]interface{}{
		//"add1": "string",
		"pointer_node":      &JsonParseNode{},
		"none_pointer_node": JsonParseNode{},
		"map_value": map[string]interface{}{
			"name":   "michael",
			"age":    "27",
			"gender": "undefined",
		},
	}
	err := serializer.WriteAdditionalData(adlData)
	assert.Nil(t, err)
	result, err := serializer.GetSerializedContent()

	stringResult := string(result[:])
	assert.Contains(t, stringResult, "\"pointer_node\":")
	assert.Contains(t, stringResult, "\"none_pointer_node\":{}")
	assert.Contains(t, stringResult, "\"name\":\"michael\"")
	assert.True(t, IsJSON("{"+stringResult+"}"))
}

func TestWriteACollectionWithNill(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "value"
	serializer.WriteStringValue("key", &value)

	prop1Value1 := internal.NewTestEntity()
	idIntValue1 := "11"
	prop1Value1.SetId(&idIntValue1)

	collection := []absser.Parsable{nil, prop1Value1}
	err := serializer.WriteCollectionOfObjectValues("", collection)

	assert.Nil(t, err)
	result, err := serializer.GetSerializedContent()

	stringResult := string(result[:])
	assert.Contains(t, stringResult, "null,")
	assert.Contains(t, stringResult, "\"key\":\"value\",[null,{\"id\":\"11\"}]")
}

func IsJSON(str string) bool {
	var js json.RawMessage
	err := json.Unmarshal([]byte(str), &js)
	return err == nil
}

func TestEscapesNewLinesInStrings(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "value\nwith\nnew\nlines"
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, "\"key\":\"value\\nwith\\nnew\\nlines\"", string(result[:]))

	fullPayload := "{" + string(result[:]) + "}"
	var parsedResult TestStruct
	parseErr := json.Unmarshal([]byte(fullPayload), &parsedResult)
	assert.Nil(t, parseErr)
	assert.Equal(t, value, parsedResult.Key)
}

func TestPreserveSeparatorsInStrings(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := `{"foo":"bar","biz":"bang"},[1,2,3,],,`
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, `"key":"{\"foo\":\"bar\",\"biz\":\"bang\"},[1,2,3,],,"`, string(result))

	fullPayload := "{" + string(result[:]) + "}"
	var parsedResult TestStruct
	parseErr := json.Unmarshal([]byte(fullPayload), &parsedResult)
	assert.Nil(t, parseErr)
	assert.Equal(t, value, parsedResult.Key)
}

func TestEscapesBackslashesInStrings(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "value\\with\\backslashes"
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, "\"key\":\"value\\\\with\\\\backslashes\"", string(result[:]))

	fullPayload := "{" + string(result[:]) + "}"
	var parsedResult TestStruct
	parseErr := json.Unmarshal([]byte(fullPayload), &parsedResult)
	assert.Nil(t, parseErr)
	assert.Equal(t, value, parsedResult.Key)
}

func TestEscapeTabAndCarriageReturnInStrings(t *testing.T) {
	doubleB := "<html lang=\"en\" style=\"min-height:100%;\t background:#ffffff\"><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\"><meta name=\"viewport\" content=\"width=device-width\"><meta name=\"eventId\" ^Mcontent=\"aad-identity-protection-weekly-digest-report-v2\"><meta name=\"messageId\" content=\"d4db577e-fe10-4bea-8e6d-164bb1ebb039\">"
	expected := "\"<html lang=\\\"en\\\" style=\\\"min-height:100%;\\t background:#ffffff\\\"><head><meta http-equiv=\\\"Content-Type\\\" content=\\\"text/html; charset=utf-8\\\"><meta name=\\\"viewport\\\" content=\\\"width=device-width\\\"><meta name=\\\"eventId\\\" ^Mcontent=\\\"aad-identity-protection-weekly-digest-report-v2\\\"><meta name=\\\"messageId\\\" content=\\\"d4db577e-fe10-4bea-8e6d-164bb1ebb039\\\">\""
	serializer := NewJsonSerializationWriter()
	err := serializer.WriteStringValue("", &doubleB)
	assert.NoError(t, err)
	result, err := serializer.GetSerializedContent()
	assert.NoError(t, err)
	converted := string(result)
	assert.Equal(t, expected, converted)
}

// TestShortEscapeSequencesInString tests that strings containing characters
// with 2-character escape sequences according to RFC 8259 section 7 are
// properly encoded as JSON.
func TestShortEscapeSequencesInString(t *testing.T) {
	// Expected results for each test are quoted since it's a JSON string.
	table := []struct {
		input    byte
		expected []byte
	}{
		{
			input:    0x22, // " character
			expected: []byte(`"\""`),
		},
		{
			input:    0x5c, // \ character
			expected: []byte(`"\\"`),
		},
		{
			input:    0x08, // backspace character
			expected: []byte(`"\b"`),
		},
		{
			input:    0x0c, // form feed character
			expected: []byte(`"\f"`),
		},
		{
			input:    0x0a, // line feed character
			expected: []byte(`"\n"`),
		},
		{
			input:    0x0d, // carriage return character
			expected: []byte(`"\r"`),
		},
		{
			input:    0x09, // tab character
			expected: []byte(`"\t"`),
		},
	}

	for _, test := range table {
		t.Run(fmt.Sprintf("0x%02X", test.input), func(t *testing.T) {
			stringInput := string(test.input)

			serializer := NewJsonSerializationWriter()
			err := serializer.WriteStringValue("", &stringInput)
			assert.NoError(t, err)

			result, err := serializer.GetSerializedContent()
			assert.NoError(t, err)

			assert.Equal(t, test.expected, result)

			assert.True(t, json.Valid(result), "valid JSON")
		})
	}
}

// TestLongEscapeSequencesInString tests that strings containing characters
// without 2-character escape sequences according to RFC 8259 section 7 are
// properly encoded as JSON.
func TestLongEscapeSequencesInString(t *testing.T) {
	// Manually adding these expected results since the code to generate them with
	// a loop would be pretty similar to the code to generate the escape sequences
	// which could make it susceptible to similar logic errors.
	table := []struct {
		input    byte
		expected []byte
	}{
		{
			input:    0x00,
			expected: []byte(`"\u0000"`),
		},
		{
			input:    0x01,
			expected: []byte(`"\u0001"`),
		},
		{
			input:    0x02,
			expected: []byte(`"\u0002"`),
		},
		{
			input:    0x03,
			expected: []byte(`"\u0003"`),
		},
		{
			input:    0x04,
			expected: []byte(`"\u0004"`),
		},
		{
			input:    0x05,
			expected: []byte(`"\u0005"`),
		},
		{
			input:    0x06,
			expected: []byte(`"\u0006"`),
		},
		{
			input:    0x07,
			expected: []byte(`"\u0007"`),
		},
		{
			input:    0x0b,
			expected: []byte(`"\u000b"`),
		},
		{
			input:    0x0e,
			expected: []byte(`"\u000e"`),
		},
		{
			input:    0x0f,
			expected: []byte(`"\u000f"`),
		},
		{
			input:    0x10,
			expected: []byte(`"\u0010"`),
		},
		{
			input:    0x11,
			expected: []byte(`"\u0011"`),
		},
		{
			input:    0x12,
			expected: []byte(`"\u0012"`),
		},
		{
			input:    0x13,
			expected: []byte(`"\u0013"`),
		},
		{
			input:    0x14,
			expected: []byte(`"\u0014"`),
		},
		{
			input:    0x15,
			expected: []byte(`"\u0015"`),
		},
		{
			input:    0x16,
			expected: []byte(`"\u0016"`),
		},
		{
			input:    0x17,
			expected: []byte(`"\u0017"`),
		},
		{
			input:    0x18,
			expected: []byte(`"\u0018"`),
		},
		{
			input:    0x19,
			expected: []byte(`"\u0019"`),
		},
		{
			input:    0x1a,
			expected: []byte(`"\u001a"`),
		},
		{
			input:    0x1b,
			expected: []byte(`"\u001b"`),
		},
		{
			input:    0x1c,
			expected: []byte(`"\u001c"`),
		},
		{
			input:    0x1d,
			expected: []byte(`"\u001d"`),
		},
		{
			input:    0x1e,
			expected: []byte(`"\u001e"`),
		},
		{
			input:    0x1f,
			expected: []byte(`"\u001f"`),
		},
	}

	for _, test := range table {
		t.Run(fmt.Sprintf("0x%02X", test.input), func(t *testing.T) {
			stringInput := string(test.input)

			serializer := NewJsonSerializationWriter()
			err := serializer.WriteStringValue("", &stringInput)
			assert.NoError(t, err)

			result, err := serializer.GetSerializedContent()
			assert.NoError(t, err)

			assert.Equal(t, test.expected, result)

			assert.True(t, json.Valid(result), "valid JSON")
		})
	}
}

func TestWriteValuesConcurrently(t *testing.T) {
	instances := 100
	output := make([][]byte, instances)

	// Use a separate function so we can just defer close the serialization
	// writer.
	serializer := func(idx int) {
		value := int64(idx)

		serializer := NewJsonSerializationWriter()
		defer serializer.Close()

		serializer.WriteInt64Value("key", &value)

		result, err := serializer.GetSerializedContent()
		require.NoError(t, err)

		output[idx] = result
	}

	for i := 0; i < instances; i++ {
		serializer(i)
	}

	for i := 0; i < instances; i++ {
		assert.Equal(t, fmt.Sprintf("\"key\":%d", i), string(output[i]))
	}
}

func TestJsonSerializationWriter_WriteNullValue(t *testing.T) {
	serializer := NewJsonSerializationWriter()

	err := serializer.WriteNullValue("name")
	assert.NoError(t, err)
	result, err := serializer.GetSerializedContent()
	assert.NoError(t, err)
	converted := string(result)

	assert.Equal(t, "\"name\":null", converted)
}

func TestJsonSerializationWriter(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	countBefore := 0
	onBefore := func(parsable absser.Parsable) error {
		countBefore++
		return nil
	}
	err := serializer.SetOnBeforeSerialization(onBefore)
	assert.NoError(t, err)

	countAfter := 0
	onAfter := func(parsable absser.Parsable) error {
		countAfter++
		return nil
	}
	err = serializer.SetOnAfterObjectSerialization(onAfter)
	assert.NoError(t, err)

	countStart := 0
	onStart := func(absser.Parsable, absser.SerializationWriter) error {
		countStart++
		return nil
	}

	err = serializer.SetOnStartObjectSerialization(onStart)
	assert.NoError(t, err)

	assert.Equal(t, 0, countBefore)
	assert.Equal(t, 0, countAfter)
	assert.Equal(t, 0, countStart)

	test := internal.NewTestEntity()
	err = serializer.WriteObjectValue("name", test)
	assert.NoError(t, err)

	assert.Equal(t, 1, countBefore)
	assert.Equal(t, 1, countAfter)
	assert.Equal(t, 1, countStart)
}

func TestWriteUntypedJson(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	untypedTestEntity := internal.NewUntypedTestEntity()
	id := "1"
	untypedTestEntity.SetId(&id)
	title := "Title"
	untypedTestEntity.SetTitle(&title)

	locationProperties := make(map[string]absser.UntypedNodeable)

	addressProperties := make(map[string]absser.UntypedNodeable)
	addressProperties["city"] = absser.NewUntypedString("Redmond")
	addressProperties["postalCode"] = absser.NewUntypedString("98052")
	addressProperties["state"] = absser.NewUntypedString("Washington")
	addressProperties["street"] = absser.NewUntypedString("NE 36th St")

	locationProperties["address"] = absser.NewUntypedObject(addressProperties)

	coordinatesProperties := make(map[string]absser.UntypedNodeable)
	coordinatesProperties["latitude"] = absser.NewUntypedDouble(47.641942)
	coordinatesProperties["longitude"] = absser.NewUntypedDouble(-122.127222)

	locationProperties["coordinates"] = absser.NewUntypedObject(coordinatesProperties)

	locationProperties["displayName"] = absser.NewUntypedString("Microsoft Building 92")
	locationProperties["floorCount"] = absser.NewUntypedInteger(int32(50))
	locationProperties["hasReception"] = absser.NewUntypedBoolean(true)
	locationProperties["contact"] = absser.NewUntypedNull()
	location := absser.NewUntypedObject(locationProperties)
	untypedTestEntity.SetLocation(location)

	keywords := make([]absser.UntypedNodeable, 2)
	firstKeywordProperties := make(map[string]absser.UntypedNodeable)
	firstKeywordProperties["created"] = absser.NewUntypedString("2023-07-26T10:41:26Z")
	firstKeywordProperties["label"] = absser.NewUntypedString("Keyword1")
	firstKeywordProperties["termGuid"] = absser.NewUntypedString("10e9cc83-b5a4-4c8d-8dab-4ada1252dd70")
	firstKeywordProperties["wssId"] = absser.NewUntypedLong(int64(6442450941))
	keywords[0] = absser.NewUntypedObject(firstKeywordProperties)

	secondKeywordProperties := make(map[string]absser.UntypedNodeable)
	secondKeywordProperties["created"] = absser.NewUntypedString("2023-07-26T10:51:26Z")
	secondKeywordProperties["label"] = absser.NewUntypedString("Keyword2")
	secondKeywordProperties["termGuid"] = absser.NewUntypedString("2cae6c6a-9bb8-4a78-afff-81b88e735fef")
	secondKeywordProperties["wssId"] = absser.NewUntypedLong(int64(6442450942))
	keywords[1] = absser.NewUntypedObject(secondKeywordProperties)

	untypedKeywordsArray := absser.NewUntypedArray(keywords)
	untypedTestEntity.SetKeywords(untypedKeywordsArray)

	extraProperties := make(map[string]absser.UntypedNodeable)
	extraProperties["createdDateTime"] = absser.NewUntypedString("2024-01-15T00:00:00+00:00")
	extra := absser.NewUntypedObject(extraProperties)
	additionalData := make(map[string]interface{})
	additionalData["extra"] = extra
	untypedTestEntity.SetAdditionalData(additionalData)

	err := serializer.WriteObjectValue("", untypedTestEntity)
	assert.NoError(t, err)
	result, err := serializer.GetSerializedContent()
	assert.NoError(t, err)
	resultString := string(result[:])
	assert.Contains(t, resultString, "\"id\":\"1\",")
	assert.Contains(t, resultString, "\"title\":\"Title\",")
	assert.Contains(t, resultString, "\"extra\":{\"createdDateTime\":\"2024-01-15T00:00:00+00:00\"}}")
	assert.Contains(t, resultString, "\"hasReception\":true")
	assert.Contains(t, resultString, "\"keywords\":[")
	assert.Contains(t, resultString, "\"wssId\":6442450942")
	assert.Contains(t, resultString, "\"hasReception\":true")
	assert.Contains(t, resultString, "\"floorCount\":50")
	assert.Contains(t, resultString, "\"termGuid\":\"10e9cc83-b5a4-4c8d-8dab-4ada1252dd70\"")

}

type TestStruct struct {
	Key string `json:"key"`
}
