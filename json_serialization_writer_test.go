package jsonserialization

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"

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

//  ByteArray values are encoded to Base64 when stored
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

func TestBufferClose(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "W/\"CQAAABYAAAAs+XSiyjZdS4Rhtwk0v1pGAAC5bsJ2\""
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.True(t, len(result) > 0)
	serializer.Close()
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

type TestStruct struct {
	Key string `json:"key"`
}
