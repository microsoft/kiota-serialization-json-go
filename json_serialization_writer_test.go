package jsonserialization

import (
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
	type Invalid string
	var value Invalid = "value"
	adlData := map[string]interface{}{
		"key": value,
	}
	err := serializer.WriteAdditionalData(adlData)
	expErr := fmt.Sprintf("unsupported AdditionalData type: %T", value)
	assert.EqualErrorf(t, err, expErr, "Error should be: %v, got: %v", expErr, err)
}

func TestEscapesNewLinesInStrings(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	value := "value\nwith\nnew\nlines"
	serializer.WriteStringValue("key", &value)
	result, err := serializer.GetSerializedContent()
	assert.Nil(t, err)
	assert.Equal(t, "\"key\":\"value\\nwith\\nnew\\nlines\"", string(result[:]))
}

func TestAdditionalDataWithEmbeddedMaps(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	data := "a string"
	aNum := int32(4)
	aBool := false
	adlData := map[string]interface{}{
		"type":        "map test",
		"isDelegated": true,
		"location": map[string]*string{"displayName": (*string)(&data),
			"locationType": (*string)(&data),
			"uniqueIdType": (*string)(&data)},
		"startDateTime": map[string]*int32{"dateTime": (*int32)(&aNum),
			"timeZone": (*int32)(&aNum)},
		"endDateTime": map[string]int32{"dateTime": *(*int32)(&aNum),
			"timeZone": *(*int32)(&aNum)},
		"meetingMessageType": &data,
		"meetingRequestType": &aNum,
		"@odata.type":        int32(13),
		"@odata.etag":        "W/\"CwAAABYAAADSEBNbUIB9RL6ePDeF3FIYAAAAAAsl\"",
		"isOutOfDate":        &aBool,
		"responseRequested":  aBool,
	}
	err := serializer.WriteAdditionalData(adlData)
	assert.NoError(t, err)
	bytes, err := serializer.GetSerializedContent()
	assert.NoError(t, err)
	assert.Greater(t, len(bytes), 0)
}

func TestAdditionalDataEmbeddedComplex(t *testing.T) {
	serializer := NewJsonSerializationWriter()
	time, err := NewJsonParseNode([]byte("\"receivedDateTime\": \"2021-10-14T09:19:02Z\","))
	if err != nil {
		t.Error("setup failure: ParseNode creation")
	}
	tz, err := NewJsonParseNode([]byte("\"timeZone\": \"UTC\""))
	if err != nil {
		t.Error("setup failure: ParseNode creation")
	}
	aBool := false
	aString := "singleInstance"
	adlData := map[string]interface{}{
		"@odata.context":     "f435c656-f8b2-4dmessages/entity",
		"type":               &aString,
		"endDateTime":        map[string]*JsonParseNode{"dateTime": time, "timeZone": tz},
		"isDelegated":        &aBool,
		"meetingMessageType": aBool,
		"responseType":       "accepted",
		"@odata.type":        "#microsoft.graph.eventMessageResponse",
		"@odata.etag":        "W/\"DAAAABYRpEYIUq+AAAfar4a\"",
		"isAllDay":           &aBool,
		"startDateTime":      map[string]*JsonParseNode{"dateTime": time, "timeZone": tz},
		"isOutOfDate":        aBool,
	}
	err = serializer.WriteAdditionalData(adlData)
	assert.NoError(t, err)
	bytes, err := serializer.GetSerializedContent()
	assert.NoError(t, err)
	assert.Greater(t, len(bytes), 0)
}
