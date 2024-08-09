package jsonserialization

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCompatible(t *testing.T) {
	cases := []struct {
		Title     string
		InputVal  interface{}
		InputType reflect.Type
		Expected  interface{}
	}{
		{
			Title:     "Valid",
			InputVal:  int(5),
			InputType: reflect.TypeOf(int8(0)),
			Expected:  true,
		},
		{
			Title:     "Incompatible",
			InputVal:  "I am a string",
			InputType: reflect.TypeOf(int8(0)),
			Expected:  false,
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			isComp := isCompatible(test.InputVal, test.InputType)
			assert.Equal(t, test.Expected, isComp)
		})
	}
}

func TestAs(t *testing.T) {
	cases := []struct {
		Title    string
		InputVal []interface{}
		Expected interface{}
		Error    error
	}{
		{
			Title: "Number",
			InputVal: []interface{}{
				int8(1),
				int16(0),
			},
			Expected: int16(1),
			Error:    nil,
		},
		{
			Title: "Interface",
			InputVal: []interface{}{
				int8(1),
				nil,
			},
			Expected: interface{}(int8(1)),
			Error:    nil,
		},
		{
			Title: "Incompatible",
			InputVal: []interface{}{
				"I am a string",
				int8(0),
			},
			Expected: int8(0),
			Error:    errors.New("value 'I am a string' is not compatible with type int8"),
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			in := test.InputVal[0]
			out := test.InputVal[1]
			err := as(in, &out)

			assert.Equal(t, test.Error, err)
			assert.Equal(t, test.Expected, out)
		})
	}
}

func TestIsNumericType(t *testing.T) {
	cases := []struct {
		Title    string
		InputVal interface{}
		Expected bool
	}{
		{
			Title:    "Int",
			InputVal: int(1),
			Expected: true,
		},
		{
			Title:    "int8",
			InputVal: int8(1),
			Expected: true,
		},
		{
			Title:    "uint8",
			InputVal: uint8(1),
			Expected: true,
		},
		{
			Title:    "int16",
			InputVal: int16(1),
			Expected: true,
		},
		{
			Title:    "uint16",
			InputVal: uint16(1),
			Expected: true,
		},
		{
			Title:    "int32",
			InputVal: int32(1),
			Expected: true,
		},
		{
			Title:    "uint32",
			InputVal: uint32(1),
			Expected: true,
		},
		{
			Title:    "int64",
			InputVal: int64(1),
			Expected: true,
		},
		{
			Title:    "uint64",
			InputVal: uint64(1),
			Expected: true,
		},
		{
			Title:    "float32",
			InputVal: float32(1.1),
			Expected: true,
		},
		{
			Title:    "float64",
			InputVal: float64(1.1),
			Expected: true,
		},
		{
			Title:    "string",
			InputVal: "1.1",
			Expected: false,
		},
		{
			Title:    "bool",
			InputVal: true,
			Expected: false,
		},
		{
			Title:    "Untyped Nil",
			InputVal: nil,
			Expected: false,
		},
		{
			Title:    "Typed Nil",
			InputVal: (*int)(nil),
			Expected: false,
		},
		{
			Title:    "Interface",
			InputVal: interface{}(int(1)),
			Expected: true,
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			isNumber := isNumericType(test.InputVal)

			assert.Equal(t, test.Expected, isNumber)
		})
	}
}

func TestIsCompatibleInt(t *testing.T) {
	cases := []struct {
		Title     string
		InputVal  interface{}
		InputType reflect.Type
		Expected  bool
	}{
		{
			Title:     "Valid",
			InputVal:  1,
			InputType: reflect.TypeOf(float64(0)),
			Expected:  true,
		},
		{
			Title:     "Too Big",
			InputVal:  300,
			InputType: reflect.TypeOf(int8(0)),
			Expected:  false,
		},
		{
			Title:     "String",
			InputVal:  "1",
			InputType: reflect.TypeOf(int8(0)),
			Expected:  false,
		},
		{
			Title:     "Nested Int",
			InputVal:  interface{}(int64(1)),
			InputType: reflect.TypeOf(int8(0)),
			Expected:  true,
		},
		{
			Title:     "Bool",
			InputVal:  true,
			InputType: reflect.TypeOf(int8(0)),
			Expected:  false,
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			isNumber := isCompatibleInt(test.InputVal, test.InputType)

			assert.Equal(t, test.Expected, isNumber)
		})
	}
}

func TestHasDecimal(t *testing.T) {
	cases := []struct {
		Title    string
		InputVal float64
		Expected bool
	}{
		{
			Title:    "Yes",
			InputVal: 1.000005,
			Expected: true,
		},
		{
			Title:    "No",
			InputVal: 1.0,
			Expected: false,
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			hasDecimal := hasDecimalPlace(test.InputVal)

			assert.Equal(t, test.Expected, hasDecimal)
		})
	}
}
