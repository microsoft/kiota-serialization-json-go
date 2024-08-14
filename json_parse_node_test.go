package jsonserialization

import (
	"errors"
	"reflect"
	"testing"

	"github.com/microsoft/kiota-serialization-json-go/internal"
	"github.com/stretchr/testify/require"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	source := "{\"someProp\": \"stringValue\",\"otherProp\": [1,2,3],\"objectProp\": {\"boolProp\": true}}"
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	someProp, err := parseNode.GetChildNode("someProp")
	if err != nil {
		t.Errorf("Error getting child node: %s", err.Error())
	}
	stringValue, err := someProp.GetStringValue()
	if err != nil {
		t.Errorf("Error getting string value: %s", err.Error())
	}
	if *stringValue != "stringValue" {
		t.Errorf("Expected value to be 'stringValue', got '%s'", *stringValue)
	}
	otherProp, err := parseNode.GetChildNode("otherProp")
	if err != nil {
		t.Errorf("Error getting child node: %s", err.Error())
	}
	arrayValue, err := otherProp.GetCollectionOfPrimitiveValues("int32")
	if err != nil {
		t.Errorf("Error getting array value: %s", err.Error())
	}
	if len(arrayValue) != 3 {
		t.Errorf("Expected array to have 3 elements, got %d", len(arrayValue))
	}
	if *(arrayValue[0].(*int32)) != 1 {
		t.Errorf("Expected array element 0 to be 1, got %d", arrayValue[0])
	}
	objectProp, err := parseNode.GetChildNode("objectProp")
	if err != nil {
		t.Errorf("Error getting child node: %s", err.Error())
	}
	boolProp, err := objectProp.GetChildNode("boolProp")
	if err != nil {
		t.Errorf("Error getting child node: %s", err.Error())
	}
	boolValue, err := boolProp.GetBoolValue()
	if err != nil {
		t.Errorf("Error getting boolean value: %s", err.Error())
	}
	if !*boolValue {
		t.Errorf("Expected value to be true, got false")
	}
}

func TestGetRawValue(t *testing.T) {
	source := `{
				"id": "2",
				"status": 200,
				"item": null,
				"phones": [1,2,3]
		  }`
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	someProp, err := parseNode.GetChildNode("item")
	value, err := someProp.GetRawValue()
	require.NoError(t, err)
	assert.Nil(t, value)

	someProp, err = parseNode.GetChildNode("status")
	value, err = someProp.GetRawValue()
	assert.Equal(t, float64(200), *value.(*float64))

	someProp, err = parseNode.GetChildNode("phones")
	value, err = someProp.GetRawValue()

	var expected []interface{}
	expected = append(expected, ref(float64(1)))
	expected = append(expected, ref(float64(2)))
	expected = append(expected, ref(float64(3)))

	assert.Equal(t, expected, value)
}

func TestNestedGetRawValue(t *testing.T) {
	source := `{
				"id": "2",
				"status": 200,
				"item": null,
				"phones": [1,2,3],
				"passwordCredentials": [
					{
						"endDateTime": "2023-07-11T14:18:14.946Z",
						"keyId": "f92ec133-34aa-49e9-b078-9f0b247d8059"
					}
				]
		  }`
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	someProp, err := parseNode.GetChildNode("passwordCredentials")
	value, err := someProp.GetRawValue()
	require.NoError(t, err)

	var expected []interface{}
	e := make(map[string]interface{})
	e["endDateTime"] = ref("2023-07-11T14:18:14.946Z")
	e["keyId"] = ref("f92ec133-34aa-49e9-b078-9f0b247d8059")
	expected = append(expected, e)

	assert.Equal(t, expected, value)
}

func TestValidEnumValue(t *testing.T) {
	source := `{
					"id": "acbb4e46-0aa9-11ee-be56-0242ac120002",
					"officeLocation": "Nairobi",
					"sensitivity": "normal"
				}`
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	parsable, err := parseNode.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	testEntity := parsable.(*internal.TestEntity)
	require.NoError(t, err)
	assert.Equal(t, "Nairobi", *testEntity.GetOfficeLocation())
	assert.Equal(t, internal.NORMAL_SENSITIVITY, *testEntity.GetSensitivity())
}

func TestInvalidEnumValueReturnsNil(t *testing.T) {
	source := `{
					"id": "acbb4e46-0aa9-11ee-be56-0242ac120002",
					"officeLocation": "Nairobi",
					"sensitivity": "invalid"
				}`
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	parsable, err := parseNode.GetObjectValue(internal.CreateTestEntityFromDiscriminator)
	testEntity := parsable.(*internal.TestEntity)
	require.NoError(t, err)
	assert.Equal(t, "Nairobi", *testEntity.GetOfficeLocation())
	assert.Nil(t, testEntity.GetSensitivity())
}

func TestNilValuesInCollections(t *testing.T) {
	source := `{
				"id": "2",
				"status": 200,
				"item": null,
				"phones": [1,2, null,3],
				"testEntities": [
					{
						"id": "acbb4e46-0aa9-11ee-be56-0242ac120002",
						"officeLocation": "Nairobi",
						"sensitivity": "personal"
					},
					null,
					{
						"id": "acbb4e46-0aa9-11ee-be56-0242ac120002",
						"officeLocation": "Nairobi",
						"sensitivity": "confidential"
					}
				]
		  }`
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	someProp, err := parseNode.GetChildNode("testEntities")
	require.NoError(t, err)
	value, err := someProp.GetCollectionOfObjectValues(internal.TestEntityDiscriminator)
	require.NoError(t, err)
	assert.Equal(t, "Nairobi", *(value[0].(*internal.TestEntity)).GetOfficeLocation())
	assert.Nil(t, value[1])

	phoneProp, err := parseNode.GetChildNode("phones")
	require.NoError(t, err)
	phonesValue, err := phoneProp.GetCollectionOfPrimitiveValues("int32")
	require.NoError(t, err)
	assert.Equal(t, int32(1), *(phonesValue[0].(*int32)))
	assert.Equal(t, nil, phonesValue[2])
}

func ref[T interface{}](t T) *T {
	return &t
}

func TestJsonParseNodeHonoursInterface(t *testing.T) {
	instance := &JsonParseNode{}
	assert.Implements(t, (*absser.ParseNode)(nil), instance)
}

func TestFunctional(t *testing.T) {
	sourceArray := []byte(FunctionalTestSource)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	if parseNode == nil {
		t.Errorf("Expected parse node to be non-nil")
	}
}

func TestParsingTime(t *testing.T) {
	source := `{
			"noZone": "2023-07-12T08:54:24",
			"withZone": "2023-07-12T09:54:24+03:00"
	  }`

	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}

	someProp, err := parseNode.GetChildNode("noZone")
	assert.Nil(t, err)
	time1, err := someProp.GetTimeValue()
	assert.Nil(t, err)
	assert.Equal(t, "2023-07-12 08:54:24", time1.Format("2006-01-02 15:04:05"))

	someProp2, err := parseNode.GetChildNode("withZone")
	assert.Nil(t, err)
	time2, err := someProp2.GetTimeValue()
	assert.Nil(t, err)
	assert.Equal(t, "2023-07-12 09:54:24", time2.Format("2006-01-02 15:04:05"))
}

func TestThrowErrorOfPrimitiveType(t *testing.T) {
	source := `{
				"id": "2",
				"status": 200,
				"item": null,
				"phones": [1,2,3]
		  }`
	sourceArray := []byte(source)
	parseNode, err := NewJsonParseNode(sourceArray)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}

	someProp, err := parseNode.GetChildNode("phones")
	_, err = someProp.GetCollectionOfPrimitiveValues("wrong.UUID")
	assert.Equal(t, "targetType wrong.UUID is not supported", err.Error())
}

func TestUntypedJsonObject(t *testing.T) {
	sourceJson := []byte(TestUntypedJson)
	parseNode, err := NewJsonParseNode(sourceJson)
	if err != nil {
		t.Errorf("Error creating parse node: %s", err.Error())
	}
	if parseNode == nil {
		t.Errorf("Expected parse node to be non-nil")
	}

	parsable, err := parseNode.GetObjectValue(internal.UntypedTestEntityDiscriminator)
	testEntity := parsable.(*internal.UntypedTestEntity)
	assert.Nil(t, err)
	assert.NotNil(t, testEntity)

	assert.Equal(t, "5", *testEntity.GetId())
	assert.Equal(t, "Project 101", *testEntity.GetTitle())
	assert.NotNil(t, testEntity.GetLocation())
	assert.NotNil(t, testEntity.GetKeywords())
	assert.Nil(t, testEntity.GetDetail())

	location := testEntity.GetLocation().(*absser.UntypedObject)
	assert.NotNil(t, location)
	locationProperties := location.GetValue()
	assert.NotNil(t, locationProperties)

	untypedDisplayName := locationProperties["displayName"].(*absser.UntypedString)
	assert.Equal(t, "Microsoft Building 92", *untypedDisplayName.GetValue())

	untypedAddress := locationProperties["address"].(*absser.UntypedObject)
	assert.NotNil(t, untypedAddress)

	untypedCount := locationProperties["floorCount"].(*absser.UntypedDouble)
	assert.Equal(t, float64(50), *untypedCount.GetValue())

	untypedBool := locationProperties["hasReception"].(*absser.UntypedBoolean)
	assert.Equal(t, true, *untypedBool.GetValue())

	untypedNull := locationProperties["contact"].(*absser.UntypedNull)
	assert.Equal(t, nil, untypedNull.GetValue())

	untypedArray := testEntity.GetKeywords().(*absser.UntypedArray)
	assert.NotNil(t, untypedArray)
	assert.Equal(t, 2, len(untypedArray.GetValue()))

	additionalData := testEntity.GetAdditionalData()
	assert.NotNil(t, additionalData)

	table := testEntity.GetTable().(*absser.UntypedArray)
	assert.NotNil(t, untypedArray)
	for _, row := range table.GetValue() {
		rowValue := row.(*absser.UntypedArray)
		assert.NotNil(t, rowValue)
		for _, cell := range rowValue.GetValue() {
			cellValue := cell.(*absser.UntypedDouble)
			assert.NotNil(t, cellValue)
		}
	}
}

func TestJsonGetStringValue(t *testing.T) {
	cases := []struct {
		Title    string
		Input    []byte
		Expected interface{}
		Error    error
	}{
		{
			Title:    "Valid",
			Input:    []byte(`"I am a string"`),
			Expected: "I am a string",
			Error:    nil,
		},
		{
			//Intentionally does not work, see https://github.com/microsoft/kiota-serialization-json-go/issues/142
			Title:    "Integer",
			Input:    []byte(`1`),
			Expected: (*string)(nil),
			Error:    errors.New("type '*float64' is not compatible with type string"),
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			var val any

			node, err := NewJsonParseNode(test.Input)
			assert.Nil(t, err)

			val, err = node.GetStringValue()

			assert.Equal(t, test.Error, err)
			v := reflect.ValueOf(val)
			if !v.IsNil() && v.Kind() == reflect.Ptr {
				val = v.Elem().Interface()
			}
			assert.Equal(t, test.Expected, val)
		})
	}
}

func TestJsonGetBoolValue(t *testing.T) {
	cases := []struct {
		Title    string
		Input    []byte
		Expected interface{}
		Error    error
	}{
		{
			Title:    "Valid",
			Input:    []byte(`true`),
			Expected: true,
			Error:    nil,
		},
		{
			Title:    "Integer",
			Input:    []byte(`1`),
			Expected: (*bool)(nil),
			Error:    errors.New("type '*float64' is not compatible with type bool"),
		},
		{
			Title:    "String",
			Input:    []byte(`"true"`),
			Expected: (*bool)(nil),
			Error:    errors.New("type '*string' is not compatible with type bool"),
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			var val any

			node, err := NewJsonParseNode(test.Input)
			assert.Nil(t, err)

			val, err = node.GetBoolValue()

			assert.Equal(t, test.Error, err)
			v := reflect.ValueOf(val)
			if !v.IsNil() && v.Kind() == reflect.Ptr {
				val = v.Elem().Interface()
			}
			assert.Equal(t, test.Expected, val)
		})
	}
}

func TestJsonGetInt8Value(t *testing.T) {
	cases := []struct {
		Title    string
		Input    []byte
		Expected interface{}
		Error    error
	}{
		{
			Title:    "Valid",
			Input:    []byte(`1`),
			Expected: int8(1),
			Error:    nil,
		},
		{
			Title:    "Bool",
			Input:    []byte(`true`),
			Expected: (*int8)(nil),
			Error:    errors.New("value 'true' is not compatible with type int8"),
		},
		{
			Title:    "String",
			Input:    []byte(`"1"`),
			Expected: (*int8)(nil),
			Error:    errors.New("value '1' is not compatible with type int8"),
		},
		{
			Title:    "Float",
			Input:    []byte(`1.1`),
			Expected: (*int8)(nil),
			Error:    errors.New("value '1.1' is not compatible with type int8"),
		},
		{
			Title:    "Too Big",
			Input:    []byte(`129`),
			Expected: (*int8)(nil),
			Error:    errors.New("value '129' is not compatible with type int8"),
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			var val any

			node, err := NewJsonParseNode(test.Input)
			assert.Nil(t, err)

			val, err = node.GetInt8Value()

			assert.Equal(t, test.Error, err)
			v := reflect.ValueOf(val)
			if !v.IsNil() && v.Kind() == reflect.Ptr {
				val = v.Elem().Interface()
			}
			assert.Equal(t, test.Expected, val)
		})
	}
}

func TestJsonGetByteValue(t *testing.T) {
	cases := []struct {
		Title    string
		Input    []byte
		Expected interface{}
		Error    error
	}{
		{
			Title:    "Valid",
			Input:    []byte(`1`),
			Expected: uint8(1),
			Error:    nil,
		},
		{
			Title:    "Bool",
			Input:    []byte(`true`),
			Expected: (*uint8)(nil),
			Error:    errors.New("value 'true' is not compatible with type uint8"),
		},
		{
			Title:    "Float",
			Input:    []byte(`1.1`),
			Expected: (*uint8)(nil),
			Error:    errors.New("value '1.1' is not compatible with type uint8"),
		},
		{
			Title:    "String",
			Input:    []byte(`"1"`),
			Expected: (*uint8)(nil),
			Error:    errors.New("value '1' is not compatible with type uint8"),
		},
		{
			Title:    "Too Big",
			Input:    []byte(`3.40283e+38`),
			Expected: (*uint8)(nil),
			Error:    errors.New("value '3.40283e+38' is not compatible with type uint8"),
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			var val any

			node, err := NewJsonParseNode(test.Input)
			assert.Nil(t, err)

			val, err = node.GetByteValue()

			assert.Equal(t, test.Error, err)
			v := reflect.ValueOf(val)
			if !v.IsNil() && v.Kind() == reflect.Ptr {
				val = v.Elem().Interface()
			}
			assert.Equal(t, test.Expected, val)
		})
	}
}

func TestJsonGetFloat32Value(t *testing.T) {
	cases := []struct {
		Title    string
		Input    []byte
		Expected interface{}
		Error    error
	}{
		{
			Title:    "Valid",
			Input:    []byte(`1`),
			Expected: float32(1),
			Error:    nil,
		},
		{
			Title:    "Bool",
			Input:    []byte(`true`),
			Expected: (*float32)(nil),
			Error:    errors.New("value 'true' is not compatible with type float32"),
		},
		{
			Title:    "String",
			Input:    []byte(`"1"`),
			Expected: (*float32)(nil),
			Error:    errors.New("value '1' is not compatible with type float32"),
		},
		{
			Title:    "Too Big",
			Input:    []byte(`3.40283e+38`),
			Expected: (*float32)(nil),
			Error:    errors.New("value '3.40283e+38' is not compatible with type float32"),
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			var val any

			node, err := NewJsonParseNode(test.Input)
			assert.Nil(t, err)

			val, err = node.GetFloat32Value()

			assert.Equal(t, test.Error, err)
			v := reflect.ValueOf(val)
			if !v.IsNil() && v.Kind() == reflect.Ptr {
				val = v.Elem().Interface()
			}
			assert.Equal(t, test.Expected, val)
		})
	}
}

func TestJsonGetFloat64Value(t *testing.T) {
	cases := []struct {
		Title    string
		Input    []byte
		Expected interface{}
		Error    error
	}{
		{
			Title:    "Valid",
			Input:    []byte(`1`),
			Expected: float64(1),
			Error:    nil,
		},
		{
			Title:    "Bool",
			Input:    []byte(`true`),
			Expected: (*float64)(nil),
			Error:    errors.New("value 'true' is not compatible with type float64"),
		},
		{
			Title:    "String",
			Input:    []byte(`"1"`),
			Expected: (*float64)(nil),
			Error:    errors.New("value '1' is not compatible with type float64"),
		},
		//NOTE: no point in checking too big, the STD JSON encoder will error out first :)
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			var val any

			node, err := NewJsonParseNode(test.Input)
			assert.Nil(t, err)

			val, err = node.GetFloat64Value()

			assert.Equal(t, test.Error, err)
			v := reflect.ValueOf(val)
			if !v.IsNil() && v.Kind() == reflect.Ptr {
				val = v.Elem().Interface()
			}
			assert.Equal(t, test.Expected, val)
		})
	}
}

func TestJsonGetInt32Value(t *testing.T) {
	cases := []struct {
		Title    string
		Input    []byte
		Expected interface{}
		Error    error
	}{
		{
			Title:    "Valid",
			Input:    []byte(`1`),
			Expected: int32(1),
			Error:    nil,
		},
		{
			Title:    "Bool",
			Input:    []byte(`true`),
			Expected: (*int32)(nil),
			Error:    errors.New("value 'true' is not compatible with type int32"),
		},
		{
			Title:    "Float",
			Input:    []byte(`1.1`),
			Expected: (*int32)(nil),
			Error:    errors.New("value '1.1' is not compatible with type int32"),
		},
		{
			Title:    "String",
			Input:    []byte(`"1"`),
			Expected: (*int32)(nil),
			Error:    errors.New("value '1' is not compatible with type int32"),
		},
		{
			Title:    "Too Big",
			Input:    []byte(`3.40283e+38`),
			Expected: (*int32)(nil),
			Error:    errors.New("value '3.40283e+38' is not compatible with type int32"),
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			var val any

			node, err := NewJsonParseNode(test.Input)
			assert.Nil(t, err)

			val, err = node.GetInt32Value()

			assert.Equal(t, test.Error, err)
			v := reflect.ValueOf(val)
			if !v.IsNil() && v.Kind() == reflect.Ptr {
				val = v.Elem().Interface()
			}
			assert.Equal(t, test.Expected, val)
		})
	}
}

func TestJsonGetInt64Value(t *testing.T) {
	cases := []struct {
		Title    string
		Input    []byte
		Expected interface{}
		Error    error
	}{
		{
			Title:    "Valid",
			Input:    []byte(`1`),
			Expected: int64(1),
			Error:    nil,
		},
		{
			Title:    "Bool",
			Input:    []byte(`true`),
			Expected: (*int64)(nil),
			Error:    errors.New("value 'true' is not compatible with type int64"),
		},
		{
			Title:    "Float",
			Input:    []byte(`1.1`),
			Expected: (*int64)(nil),
			Error:    errors.New("value '1.1' is not compatible with type int64"),
		},
		{
			Title:    "Too Big",
			Input:    []byte(`3.40283e+38`),
			Expected: (*int64)(nil),
			Error:    errors.New("value '3.40283e+38' is not compatible with type int64"),
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			var val any

			node, err := NewJsonParseNode(test.Input)
			assert.Nil(t, err)

			val, err = node.GetInt64Value()

			assert.Equal(t, test.Error, err)
			v := reflect.ValueOf(val)
			if !v.IsNil() && v.Kind() == reflect.Ptr {
				val = v.Elem().Interface()
			}
			assert.Equal(t, test.Expected, val)
		})
	}
}

const TestUntypedJson = "{\r\n" +
	"    \"@odata.context\": \"https://graph.microsoft.com/v1.0/$metadata#sites('contoso.sharepoint.com')/lists('fa631c4d-ac9f-4884-a7f5-13c659d177e3')/items('1')/fields/$entity\",\r\n" +
	"    \"id\": \"5\",\r\n" +
	"    \"title\": \"Project 101\",\r\n" +
	"    \"location\": {\r\n" +
	"        \"address\": {\r\n" +
	"            \"city\": \"Redmond\",\r\n" +
	"            \"postalCode\": \"98052\",\r\n" +
	"            \"state\": \"Washington\",\r\n" +
	"            \"street\": \"NE 36th St\"\r\n" +
	"        },\r\n" +
	"        \"coordinates\": {\r\n" +
	"            \"latitude\": 47.641942,\r\n" +
	"            \"longitude\": -122.127222\r\n" +
	"        },\r\n" +
	"        \"displayName\": \"Microsoft Building 92\",\r\n" +
	"        \"floorCount\": 50,\r\n" +
	"        \"hasReception\": true,\r\n" +
	"        \"contact\": null\r\n" +
	"    },\r\n" +
	"    \"keywords\": [\r\n" +
	"        {\r\n" +
	"            \"created\": \"2023-07-26T10:41:26Z\",\r\n" +
	"            \"label\": \"Keyword1\",\r\n" +
	"            \"termGuid\": \"10e9cc83-b5a4-4c8d-8dab-4ada1252dd70\",\r\n" +
	"            \"wssId\": 6442450942\r\n" +
	"        },\r\n" +
	"        {\r\n" +
	"            \"created\": \"2023-07-26T10:51:26Z\",\r\n" +
	"            \"label\": \"Keyword2\",\r\n" +
	"            \"termGuid\": \"2cae6c6a-9bb8-4a78-afff-81b88e735fef\",\r\n" +
	"            \"wssId\": 6442450943\r\n" +
	"        }\r\n" +
	"    ],\r\n" +
	"    \"detail\": null,\r\n" +
	"    \"table\": [[1,2,3],[4,5,6],[7,8,9]],\r\n" +
	"    \"extra\": {\r\n" +
	"        \"createdDateTime\":\"2024-01-15T00:00:00\\u002B00:00\"\r\n" +
	"    }\r\n" +
	"}"

const FunctionalTestSource = "{" +
	"\"@odata.context\": \"https://graph.microsoft.com/v1.0/$metadata#users('vincent%40biret365.onmicrosoft.com')/messages\"," +
	"\"@odata.nextLink\": \"https://graph.microsoft.com/v1.0/users/vincent@biret365.onmicrosoft.com/messages?$skip=10\"," +
	"\"value\": [" +
	"{" +
	"\"@odata.etag\": \"W/\\\"CQAAABYAAAAs+XSiyjZdS4Rhtwk0v1pGAAA4Xv0v\\\"\"," +
	"\"id\": \"AAMkAGNmMGZiNjM5LTZmMDgtNGU2OS1iYmUwLWYwZDc4M2ZkOGY1ZQBGAAAAAAAK20ulGawAT7z-yx90ohp-BwAs_XSiyjZdS4Rhtwk0v1pGAAAAAAEMAAAs_XSiyjZdS4Rhtwk0v1pGAAA4dw6TAAA=\"," +
	"\"createdDateTime\": \"2021-10-14T09:19:01Z\"," +
	"\"lastModifiedDateTime\": \"2021-10-14T09:19:03Z\"," +
	"\"changeKey\": \"CQAAABYAAAAs+XSiyjZdS4Rhtwk0v1pGAAA4Xv0v\"," +
	"\"categories\": []," +
	"\"receivedDateTime\": \"2021-10-14T09:19:02Z\"," +
	"\"sentDateTime\": \"2021-10-14T09:18:59Z\"," +
	"\"hasAttachments\": false," +
	"\"internetMessageId\": \"<608fed24166f421aa1e27a6c822074ba-JFBVALKQOJXWILKNK4YVA7CPGM3DKTLFONZWCZ3FINSW45DFOJ6E2ZLTONQWOZKDMVXHIZLSL5GUGMRZGEYDQOD4KNWXI4A=@microsoft.com>\"," +
	"\"subject\": \"Major update from Message center\"," +
	"\"bodyPreview\": \"(Updated) Microsoft 365 Compliance Center Core eDiscovery - Search by ID list retirementMC291088 · BIRET365Updated October 13, 2021: We have updated this message with additional details for clarity.We will be retiring the option to Search by ID,\"," +
	"\"importance\": \"normal\"," +
	"\"parentFolderId\": \"AQMkAGNmMGZiNjM5LTZmMDgtNGU2OS1iYgBlMC1mMGQ3ODNmZDhmNWUALgAAAwrbS6UZrABPvP-LH3SiGn8BACz5dKLKNl1LhGG3CTS-WkYAAAIBDAAAAA==\"," +
	"\"conversationId\": \"AAQkAGNmMGZiNjM5LTZmMDgtNGU2OS1iYmUwLWYwZDc4M2ZkOGY1ZQAQANari86tqeZDsqpmA19AXLQ=\"," +
	"\"conversationIndex\": \"AQHXwNyG1quLzq2p5kOyqmYDX0BctA==\"," +
	"\"isDeliveryReceiptRequested\": null," +
	"\"isReadReceiptRequested\": false," +
	"\"isRead\": false," +
	"\"isDraft\": false," +
	"\"webLink\": \"https://outlook.office365.com/owa/?ItemID=AAMkAGNmMGZiNjM5LTZmMDgtNGU2OS1iYmUwLWYwZDc4M2ZkOGY1ZQBGAAAAAAAK20ulGawAT7z%2Fyx90ohp%2FBwAs%2BXSiyjZdS4Rhtwk0v1pGAAAAAAEMAAAs%2BXSiyjZdS4Rhtwk0v1pGAAA4dw6TAAA%3D&exvsurl=1&viewmodel=ReadMessageItem\"," +
	"\"inferenceClassification\": \"other\"," +
	"\"body\": {" +
	"\"contentType\": \"html\"," +
	"\"content\": \"<html><head><meta http-equiv=\\\"Content-Type\\\" content=\\\"text/html; charset=utf-8\\\"><meta name=\\\"viewport\\\" content=\\\"width=device-width, initial-scale=1\\\"><meta content=\\\"IE=edge\\\"><style><!--body, table, td{font-family:Segoe UI,Helvetica,Arial,sans-serif!important}a{color:#006CBE;text-decoration:none}--></style></head><body><div style=\\\"background:white; min-height:100vh; color:#323130; font-size:14px\\\"><table border=\\\"0\\\" cellpadding=\\\"0\\\" cellspacing=\\\"0\\\" width=\\\"100%\\\" height=\\\"100%\\\"><tbody><tr><td></td><td width=\\\"640\\\"><table border=\\\"0\\\" cellpadding=\\\"0\\\" cellspacing=\\\"0\\\" style=\\\"min-width:100%; background:white\\\"><tbody><tr><td style=\\\"padding:24px 24px 45px\\\"><img src=\\\"https://eus-contentstorage.osi.office.net/images/retailer.images/centralizeddeployment/logos/112fec798b78aa02.png\\\" width=\\\"100\\\" height=\\\"21\\\" alt=\\\"Microsoft\\\"> </td></tr><tr><td style=\\\"font-size:28px; padding:0 24px; font-weight:bold; color:#000000\\\">(Updated) Microsoft 365 Compliance Center Core eDiscovery - Search by ID list retirement</td></tr><tr><td style=\\\"color:#323130; padding:20px 24px 40px 24px\\\"><span style=\\\"font-weight:600\\\">MC291088 · BIRET365</span></td></tr><tr><td style=\\\"padding:0 24px 44px\\\"><div><p style=\\\"margin-top:0\\\">Updated October 13, 2021: We have updated this message with additional details for clarity.</p><p>We will be retiring the option to Search by ID list, as it is not functioning to an adequate level and creates significant challenges for organizations who depend on consistent and repeatable results for eDiscovery workflows.<br></p><p><b style=\\\"font-weight:600\\\">When will this happen:</b></p><p>We will begin making this change in mid-November and expect to complete by the end of November.</p><p><b style=\\\"font-weight:600\\\">How this will affect your organization:</b><br></p><p>You are receiving this message because our reporting indicates your organization may be using Search by ID list.</p><p>Once this change is made, the option to Search by ID list will be removed. We suggest focusing on search by query, condition and/or locations rather that ID.</p><p><b style=\\\"font-weight:600\\\">What you need to do to prepare:</b><br></p><p>To fix this problem you need to review your eDiscovery search process, and update the workflow to focus on search by Subjects and dates rather than Search by ID list. Upon export from Core eDiscovery you can explore options to refine to only the messages of interest. </p><p></p><p>Click Additional Information to find out more.<br></p><a href=\\\"https://docs.microsoft.com/microsoft-365/compliance/search-for-content-in-core-ediscovery?view=o365-worldwide\\\" title=\\\"Additional Information\\\">Additional Information</a> </div><div style=\\\"padding-top:3px\\\"><a href=\\\"https://admin.microsoft.com/AdminPortal/home#/MessageCenter/:/messages/MC291088?MCLinkSource=MajorUpdate\\\" title=\\\"view message\\\" target=\\\"_blank\\\">View this message in the Microsoft 365 admin center</a> </div></td></tr><tr><td><table border=\\\"0\\\" cellpadding=\\\"0\\\" cellspacing=\\\"0\\\" width=\\\"100%\\\" style=\\\"min-width:100%; background-color:#F3F2F1\\\"><tbody><tr><td style=\\\"padding:44px 24px 3px; font-size:10px; color:#484644\\\">You're subscribed to this email using vincent@biret365.onmicrosoft.com. If you're an IT admin, you're subscribed by default, but you can <a href=\\\"https://admin.microsoft.com/adminportal/home#/MessageCenter/:/mcpreferences\\\" target=\\\"_blank\\\">unsubscribe at any time</a>. If you're not an IT admin, ask your admin to remove your email address from Microsoft 365 message center preferences.<br><br><a href=\\\"https://docs.microsoft.com/en-us/microsoft-365/admin/manage/language-translation-for-message-center-posts?view=o365-worldwide\\\" target=\\\"_blank\\\">How to view translated messages</a><br></td></tr><tr><td style=\\\"padding:25px 24px 24px; font-size:12px\\\"><div style=\\\"color:#696969\\\">This is a mandatory service communication. To set your contact preferences or to unsubcribe from other communications, visit the <a href=\\\"https://go.microsoft.com/fwlink/?LinkId=243189\\\" target=\\\"_blank\\\" style=\\\"color:#696969; text-decoration:underline; text-decoration-color:#696969\\\">Promotional Communications Manager</a>. <a href=\\\"https://go.microsoft.com/fwlink/?LinkId=521839\\\" target=\\\"_blank\\\" style=\\\"color:#696969; text-decoration:underline; text-decoration-color:#696969\\\">Privacy statement</a>. <br><br>Il s’agit de communications obligatoires. Pour configurer vos préférences de contact pour d’autres communications, accédez au <a href=\\\"https://go.microsoft.com/fwlink/?LinkId=243189\\\" target=\\\"_blank\\\" style=\\\"color:#696969; text-decoration:underline; text-decoration-color:#696969\\\">gestionnaire de communications promotionnelles</a>. <a href=\\\"https://go.microsoft.com/fwlink/?LinkId=521839\\\" target=\\\"_blank\\\" style=\\\"color:#696969; text-decoration:underline; text-decoration-color:#696969\\\">Déclaration de confidentialité</a>. </div><div style=\\\"color:#696969; margin-top:10px; margin-bottom:13px\\\">Microsoft Corporation, One Microsoft Way, Redmond WA 98052 USA</div><img src=\\\"https://eus-contentstorage.osi.office.net/images/retailer.images/centralizeddeployment/logos/112fec798b78aa02.png\\\" width=\\\"94\\\" height=\\\"20\\\" alt=\\\"Microsoft\\\"> </td></tr></tbody></table></td></tr></tbody></table></td><td></td></tr></tbody></table></div><img src=\\\"https://mucp.api.account.microsoft.com/m/v2/v?d=AIAAD2ON6I4P6T45JIHQXRZ6AI7WMQVRDMGBPOFLIPLXZDLYEKNQK44CEBYSPPTPDHET337ASHWG3BMEXD6NQZGTF442DPYPANRAMYRCB5XW3VUZYYL7MXCMJU7NIFJFML3F22PJFGPVVKXDWKRH374HXHZFHRY&amp;i=AIAADOZFMOPSOOEFOUHZD4HWEDARG3W3DMLBKJLS4RUJB6O5L7UJYE5NWIJQFRZTMSB74FMTRBBXRGSZEHD6UYCOLJNM7JTG27THR2WYKQWVGJXJGXJDIRHKWQDFKHWPZPZGXDKOGME5EPT3MJK3LLV7VUODVXG2VLJW5SS6POXQKSQXJWFFBHDP6VMQQEX6MHHWYLSUJG4EPHC4U23LQ7P2IKBLOLB5TTYXB5WQPHDYUDO6WN7BVWK4JGZFE7JOGWQTWAGYP7NKV7L3W3XV2W2E7NOXLUQ\\\" width=\\\"1\\\" height=\\\"1\\\" tabindex=\\\"-1\\\" aria-hidden=\\\"true\\\" alt=\\\"\\\"> </body></html>\"" +
	"}," +
	"\"sender\": {" +
	"\"emailAddress\": {" +
	"\"name\": \"Microsoft 365 Message center\"," +
	"\"address\": \"o365mc@microsoft.com\"" +
	"}" +
	"}," +
	"\"from\": {" +
	"\"emailAddress\": {" +
	"\"name\": \"Microsoft 365 Message center\"," +
	"\"address\": \"o365mc@microsoft.com\"" +
	"}" +
	"}," +
	"\"toRecipients\": [" +
	"{" +
	"\"emailAddress\": {" +
	"\"name\": \"Vincent BIRET\"," +
	"\"address\": \"vincent@biret365.onmicrosoft.com\"" +
	"}" +
	"}" +
	"]," +
	"\"ccRecipients\": []," +
	"\"bccRecipients\": []," +
	"\"replyTo\": []," +
	"\"flag\": {" +
	"\"flagStatus\": \"notFlagged\"" +
	"}" +
	"}" +
	"]" +
	"}"
