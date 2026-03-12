// Package jsonserialization is the default Kiota serialization implementation for JSON.
// It relies on the standard Go JSON library.
package jsonserialization

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// JsonParseNode is a ParseNode implementation for JSON.
type JsonParseNode struct {
	value                     interface{}
	onBeforeAssignFieldValues absser.ParsableAction
	onAfterAssignFieldValues  absser.ParsableAction
}

// tokenToValue converts a JSON token to either a raw primitive value (to avoid JsonParseNode
// allocation for primitives) or a *JsonParseNode for complex types (objects and arrays).
// This is used when building parse trees to reduce allocations.
func tokenToValue(decoder *json.Decoder, token json.Token) (interface{}, error) {
	switch t := token.(type) {
	case json.Delim:
		node, err := loadJsonTreeFromToken(decoder, t)
		return node, err
	case float64:
		f := t
		return &f, nil
	case string:
		s := t
		return &s, nil
	case bool:
		b := t
		return &b, nil
	case json.Number:
		i, err := t.Int64()
		if err == nil {
			return &i, nil
		}
		f, err := t.Float64()
		if err == nil {
			return &f, nil
		}
		return nil, errors.New("failed to parse number token")
	case int8:
		v := t
		return &v, nil
	case byte:
		v := t
		return &v, nil
	case float32:
		v := t
		return &v, nil
	case int32:
		v := t
		return &v, nil
	case int64:
		v := t
		return &v, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown token type during parsing: %T", token)
	}
}

// NewJsonParseNode creates a new JsonParseNode.
func NewJsonParseNode(content []byte) (*JsonParseNode, error) {
	if len(content) == 0 {
		return nil, errors.New("content is empty")
	}
	if !json.Valid(content) {
		return nil, errors.New("invalid json type")
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	value, err := loadJsonTree(decoder)
	return value, err
}

func loadJsonTree(decoder *json.Decoder) (*JsonParseNode, error) {
	token, err := decoder.Token()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return loadJsonTreeFromToken(decoder, token)
}

// loadJsonTreeFromToken builds a JsonParseNode from an already-consumed token.
// For object and array delimiters, it reads the remaining tokens from the decoder.
// For primitive tokens it wraps the value directly in a JsonParseNode.
// Primitive values inside objects and arrays are stored as raw values (not wrapped
// in *JsonParseNode) to reduce allocations.
func loadJsonTreeFromToken(decoder *json.Decoder, token json.Token) (*JsonParseNode, error) {
	switch t := token.(type) {
	case json.Delim:
		switch t {
		case '{':
			v := make(map[string]interface{})
			for decoder.More() {
				key, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				keyStr, ok := key.(string)
				if !ok {
					return nil, errors.New("key is not a string")
				}
				valToken, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				childValue, err := tokenToValue(decoder, valToken)
				if err != nil {
					return nil, err
				}
				v[keyStr] = childValue
			}
			decoder.Token() // skip the closing curly
			return &JsonParseNode{value: v}, nil
		case '[':
			v := make([]interface{}, 0)
			for decoder.More() {
				elemToken, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				elem, err := tokenToValue(decoder, elemToken)
				if err != nil {
					return nil, err
				}
				v = append(v, elem)
			}
			decoder.Token() // skip the closing bracket
			return &JsonParseNode{value: v}, nil
		}
	case json.Number:
		i, err := t.Int64()
		if err == nil {
			return &JsonParseNode{value: &i}, nil
		}
		f, err := t.Float64()
		if err == nil {
			return &JsonParseNode{value: &f}, nil
		}
		return nil, errors.New("failed to parse number")
	case string:
		s := t
		return &JsonParseNode{value: &s}, nil
	case bool:
		b := t
		return &JsonParseNode{value: &b}, nil
	case int8:
		v := t
		return &JsonParseNode{value: &v}, nil
	case byte:
		v := t
		return &JsonParseNode{value: &v}, nil
	case float64:
		f := t
		return &JsonParseNode{value: &f}, nil
	case float32:
		v := t
		return &JsonParseNode{value: &v}, nil
	case int32:
		v := t
		return &JsonParseNode{value: &v}, nil
	case int64:
		v := t
		return &JsonParseNode{value: &v}, nil
	case nil:
		return nil, nil
	}
	return nil, nil
}

// SetValue is obsolete, parse nodes are not meant to be settable externally
func (n *JsonParseNode) SetValue(value interface{}) {
	n.setValue(value)
}

// setValue sets the value represented by the node
func (n *JsonParseNode) setValue(value interface{}) {
	n.value = value
}

// GetChildNode returns a new parse node for the given identifier.
func (n *JsonParseNode) GetChildNode(index string) (absser.ParseNode, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if index == "" {
		return nil, errors.New("index is empty")
	}
	childNodes, ok := n.value.(map[string]interface{})
	if !ok || len(childNodes) == 0 {
		return nil, nil
	}

	rawChild, exists := childNodes[index]
	if !exists {
		return nil, nil
	}

	var childNode *JsonParseNode
	if rawChild == nil {
		// JSON null value: return a typed nil so callers can still invoke methods on the
		// returned ParseNode interface without a nil-interface panic.
		return (*JsonParseNode)(nil), nil
	} else if jn, ok := rawChild.(*JsonParseNode); ok {
		childNode = jn
	} else {
		// Raw primitive value – wrap on demand to avoid pre-allocation
		childNode = &JsonParseNode{value: rawChild}
	}

	if childNode != nil {
		err := childNode.SetOnBeforeAssignFieldValues(n.GetOnBeforeAssignFieldValues())
		if err != nil {
			return nil, err
		}
		err = childNode.SetOnAfterAssignFieldValues(n.GetOnAfterAssignFieldValues())
		if err != nil {
			return nil, err
		}
	}

	return childNode, nil
}

// GetObjectValue returns the Parsable value from the node.
func (n *JsonParseNode) GetObjectValue(ctor absser.ParsableFactory) (absser.Parsable, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if ctor == nil {
		return nil, errors.New("constructor is nil")
	}
	result, err := ctor(n)
	if err != nil {
		return nil, err
	}

	_, isUntypedNode := result.(absser.UntypedNodeable)
	if isUntypedNode {
		switch value := n.value.(type) {
		case *bool:
			return absser.NewUntypedBoolean(*value), nil
		case *string:
			return absser.NewUntypedString(*value), nil
		case *float32:
			return absser.NewUntypedFloat(*value), nil
		case *float64:
			return absser.NewUntypedDouble(*value), nil
		case *int32:
			return absser.NewUntypedInteger(*value), nil
		case *int64:
			return absser.NewUntypedLong(*value), nil
		case nil:
			return absser.NewUntypedNull(), nil
		case map[string]interface{}:
			properties := make(map[string]absser.UntypedNodeable)
			for key, rawVal := range value {
				var parsable absser.Parsable
				if rawVal == nil {
					parsable = absser.NewUntypedNull()
				} else if jn, ok := rawVal.(*JsonParseNode); ok {
					parsable, err = jn.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
					if err != nil {
						return nil, errors.New("cannot parse object value")
					}
					if parsable == nil {
						parsable = absser.NewUntypedNull()
					}
				} else {
					parsable = rawToUntypedNodeable(rawVal)
				}
				if property, ok := parsable.(absser.UntypedNodeable); ok {
					properties[key] = property
				}
			}
			return absser.NewUntypedObject(properties), nil
		case []interface{}:
			collection := make([]absser.UntypedNodeable, len(value))
			for index, rawElem := range value {
				var parsable absser.Parsable
				if rawElem == nil {
					parsable = absser.NewUntypedNull()
				} else if jn, ok := rawElem.(*JsonParseNode); ok {
					parsable, err = jn.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
					if err != nil {
						return nil, errors.New("cannot parse object value")
					}
					if parsable == nil {
						parsable = absser.NewUntypedNull()
					}
				} else {
					parsable = rawToUntypedNodeable(rawElem)
				}
				if property, ok := parsable.(absser.UntypedNodeable); ok {
					collection[index] = property
				}
			}
			return absser.NewUntypedArray(collection), nil
		default:
			return absser.NewUntypedNode(value), nil
		}
	}

	abstractions.InvokeParsableAction(n.GetOnBeforeAssignFieldValues(), result)
	properties, ok := n.value.(map[string]interface{})
	fields := result.GetFieldDeserializers()
	if ok && len(properties) != 0 {
		itemAsHolder, isHolder := result.(absser.AdditionalDataHolder)
		var itemAdditionalData map[string]interface{}
		if isHolder {
			itemAdditionalData = itemAsHolder.GetAdditionalData()
			if itemAdditionalData == nil {
				itemAdditionalData = make(map[string]interface{})
				itemAsHolder.SetAdditionalData(itemAdditionalData)
			}
		}

		for key, rawValue := range properties {
			field := fields[key]
			if field == nil {
				if rawValue != nil && isHolder {
					if jn, ok := rawValue.(*JsonParseNode); ok {
						rv, err := jn.GetRawValue()
						if err != nil {
							return nil, err
						}
						itemAdditionalData[key] = rv
					} else {
						// Raw primitive stored directly – no need to create a JsonParseNode
						itemAdditionalData[key] = rawValue
					}
				}
			} else {
				// Wrap raw values in a JsonParseNode on demand only for known fields
				var childNode *JsonParseNode
				if rawValue == nil {
					childNode = nil
				} else if jn, ok := rawValue.(*JsonParseNode); ok {
					childNode = jn
				} else {
					childNode = &JsonParseNode{value: rawValue}
				}
				if childNode != nil {
					err := childNode.SetOnBeforeAssignFieldValues(n.GetOnBeforeAssignFieldValues())
					if err != nil {
						return nil, err
					}
					err = childNode.SetOnAfterAssignFieldValues(n.GetOnAfterAssignFieldValues())
					if err != nil {
						return nil, err
					}
				}
				err := field(childNode)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	abstractions.InvokeParsableAction(n.GetOnAfterAssignFieldValues(), result)
	return result, nil
}

// GetCollectionOfObjectValues returns the collection of Parsable values from the node.
func (n *JsonParseNode) GetCollectionOfObjectValues(ctor absser.ParsableFactory) ([]absser.Parsable, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if ctor == nil {
		return nil, errors.New("ctor is nil")
	}
	nodes, ok := n.value.([]interface{})
	if !ok {
		return nil, errors.New("value is not a collection")
	}
	result := make([]absser.Parsable, len(nodes))
	for i, rawElem := range nodes {
		if rawElem == nil {
			result[i] = nil
			continue
		}
		jn, ok := rawElem.(*JsonParseNode)
		if !ok {
			return nil, errors.New("collection element is not a parse node")
		}
		val, err := jn.GetObjectValue(ctor)
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}

// GetCollectionOfPrimitiveValues returns the collection of primitive values from the node.
func (n *JsonParseNode) GetCollectionOfPrimitiveValues(targetType string) ([]interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if targetType == "" {
		return nil, errors.New("targetType is empty")
	}
	nodes, ok := n.value.([]interface{})
	if !ok {
		return nil, errors.New("value is not a collection")
	}
	result := make([]interface{}, len(nodes))
	for i, rawElem := range nodes {
		if rawElem == nil {
			result[i] = nil
			continue
		}
		if jn, ok := rawElem.(*JsonParseNode); ok {
			// Complex node (e.g. nested array/object used as a primitive collection element)
			val, err := jn.getPrimitiveValue(targetType)
			if err != nil {
				return nil, err
			}
			result[i] = val
		} else {
			// Raw primitive value stored without a JsonParseNode wrapper – convert directly
			// to avoid allocating an intermediate node.
			val, err := rawToPrimitiveValue(rawElem, targetType)
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
	}
	return result, nil
}

func (n *JsonParseNode) getPrimitiveValue(targetType string) (interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	switch targetType {
	case "string":
		return n.GetStringValue()
	case "bool":
		return n.GetBoolValue()
	case "uint8":
		return n.GetInt8Value()
	case "byte":
		return n.GetByteValue()
	case "float32":
		return n.GetFloat32Value()
	case "float64":
		return n.GetFloat64Value()
	case "int32":
		return n.GetInt32Value()
	case "int64":
		return n.GetInt64Value()
	case "time":
		return n.GetTimeValue()
	case "timeonly":
		return n.GetTimeOnlyValue()
	case "dateonly":
		return n.GetDateOnlyValue()
	case "isoduration":
		return n.GetISODurationValue()
	case "uuid":
		return n.GetUUIDValue()
	case "base64":
		return n.GetByteArrayValue()
	default:
		return nil, fmt.Errorf("targetType %s is not supported", targetType)
	}
}

// rawToPrimitiveValue converts a raw primitive value (stored without a JsonParseNode wrapper)
// to the requested target type. This avoids allocating an intermediate JsonParseNode when
// processing collections of primitive values.
func rawToPrimitiveValue(rawValue interface{}, targetType string) (interface{}, error) {
	switch targetType {
	case "string":
		var val string
		if err := as(rawValue, &val); err != nil {
			return nil, err
		}
		return &val, nil
	case "bool":
		var val bool
		if err := as(rawValue, &val); err != nil {
			return nil, err
		}
		return &val, nil
	case "uint8":
		// The "uint8" target type maps to int8 to align with GetInt8Value which returns *int8.
		// This mirrors the behaviour of getPrimitiveValue / GetInt8Value.
		var val int8
		if err := as(rawValue, &val); err != nil {
			return nil, err
		}
		return &val, nil
	case "byte":
		var val byte
		if err := as(rawValue, &val); err != nil {
			return nil, err
		}
		return &val, nil
	case "float32":
		var val float32
		if err := as(rawValue, &val); err != nil {
			return nil, err
		}
		return &val, nil
	case "float64":
		var val float64
		if err := as(rawValue, &val); err != nil {
			return nil, err
		}
		return &val, nil
	case "int32":
		var val int32
		if err := as(rawValue, &val); err != nil {
			return nil, err
		}
		return &val, nil
	case "int64":
		var val int64
		if err := as(rawValue, &val); err != nil {
			return nil, err
		}
		return &val, nil
	case "time":
		// For time/date types, the raw value should be a *string; delegate to a temporary node
		tmpNode := &JsonParseNode{value: rawValue}
		return tmpNode.GetTimeValue()
	case "timeonly":
		tmpNode := &JsonParseNode{value: rawValue}
		return tmpNode.GetTimeOnlyValue()
	case "dateonly":
		tmpNode := &JsonParseNode{value: rawValue}
		return tmpNode.GetDateOnlyValue()
	case "isoduration":
		tmpNode := &JsonParseNode{value: rawValue}
		return tmpNode.GetISODurationValue()
	case "uuid":
		tmpNode := &JsonParseNode{value: rawValue}
		return tmpNode.GetUUIDValue()
	case "base64":
		tmpNode := &JsonParseNode{value: rawValue}
		return tmpNode.GetByteArrayValue()
	default:
		return nil, fmt.Errorf("targetType %s is not supported", targetType)
	}
}

// GetCollectionOfEnumValues returns the collection of Enum values from the node.
func (n *JsonParseNode) GetCollectionOfEnumValues(parser absser.EnumFactory) ([]interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	nodes, ok := n.value.([]interface{})
	if !ok {
		return nil, errors.New("value is not a collection")
	}
	result := make([]interface{}, len(nodes))
	for i, rawElem := range nodes {
		if rawElem == nil {
			result[i] = nil
			continue
		}
		var strVal *string
		if jn, ok := rawElem.(*JsonParseNode); ok {
			var err error
			strVal, err = jn.GetStringValue()
			if err != nil {
				return nil, err
			}
		} else if sp, ok := rawElem.(*string); ok {
			strVal = sp
		} else {
			return nil, fmt.Errorf("enum collection element has unexpected type %T", rawElem)
		}
		if strVal == nil {
			result[i] = nil
			continue
		}
		val, err := parser(*strVal)
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}

// GetStringValue returns a String value from the nodes.
func (n *JsonParseNode) GetStringValue() (*string, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}

	val, ok := n.value.(*string)
	if !ok {
		return nil, fmt.Errorf("type '%T' is not compatible with type string", n.value)
	}
	return val, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *JsonParseNode) GetBoolValue() (*bool, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}

	val, ok := n.value.(*bool)
	if !ok {
		return nil, fmt.Errorf("type '%T' is not compatible with type bool", n.value)
	}
	return val, nil
}

// GetInt8Value returns a int8 value from the nodes.
func (n *JsonParseNode) GetInt8Value() (*int8, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val int8

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *JsonParseNode) GetByteValue() (*byte, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val byte

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetFloat32Value returns a Float32 value from the nodes.
func (n *JsonParseNode) GetFloat32Value() (*float32, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val float32

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetFloat64Value returns a Float64 value from the nodes.
func (n *JsonParseNode) GetFloat64Value() (*float64, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val float64

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetInt32Value returns a Int32 value from the nodes.
func (n *JsonParseNode) GetInt32Value() (*int32, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val int32

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetInt64Value returns a Int64 value from the nodes.
func (n *JsonParseNode) GetInt64Value() (*int64, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val int64

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetTimeValue returns a Time value from the nodes.
func (n *JsonParseNode) GetTimeValue() (*time.Time, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	// if string does not have timezone information, add local timezone
	if len(*v) == 19 {
		*v = *v + time.Now().Format("-07:00")
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	return &parsed, err
}

// GetISODurationValue returns a ISODuration value from the nodes.
func (n *JsonParseNode) GetISODurationValue() (*absser.ISODuration, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseISODuration(*v)
}

// GetTimeOnlyValue returns a TimeOnly value from the nodes.
func (n *JsonParseNode) GetTimeOnlyValue() (*absser.TimeOnly, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseTimeOnly(*v)
}

// GetDateOnlyValue returns a DateOnly value from the nodes.
func (n *JsonParseNode) GetDateOnlyValue() (*absser.DateOnly, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseDateOnly(*v)
}

// GetUUIDValue returns a UUID value from the nodes.
func (n *JsonParseNode) GetUUIDValue() (*uuid.UUID, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	parsed, err := uuid.Parse(*v)
	return &parsed, err
}

// GetEnumValue returns a Enum value from the nodes.
func (n *JsonParseNode) GetEnumValue(parser absser.EnumFactory) (interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	s, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return parser(*s)
}

// GetByteArrayValue returns a ByteArray value from the nodes.
func (n *JsonParseNode) GetByteArrayValue() ([]byte, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	s, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(*s)
}

// GetRawValue returns a ByteArray value from the nodes.
func (n *JsonParseNode) GetRawValue() (interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	switch v := n.value.(type) {
	case *JsonParseNode:
		return v.GetRawValue()
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, x := range v {
			if x == nil {
				result[i] = nil
				continue
			}
			if jn, ok := x.(*JsonParseNode); ok {
				val, err := jn.GetRawValue()
				if err != nil {
					return nil, err
				}
				result[i] = val
			} else {
				// Raw primitive – return as-is
				result[i] = x
			}
		}
		return result, nil
	case map[string]interface{}:
		m := make(map[string]interface{})
		for key, element := range v {
			if element == nil {
				m[key] = nil
				continue
			}
			if jn, ok := element.(*JsonParseNode); ok {
				elementVal, err := jn.GetRawValue()
				if err != nil {
					return nil, err
				}
				m[key] = elementVal
			} else {
				// Raw primitive – return as-is
				m[key] = element
			}
		}
		return m, nil
	default:
		return n.value, nil
	}
}

func (n *JsonParseNode) GetOnBeforeAssignFieldValues() absser.ParsableAction {
	return n.onBeforeAssignFieldValues
}

func (n *JsonParseNode) SetOnBeforeAssignFieldValues(action absser.ParsableAction) error {
	n.onBeforeAssignFieldValues = action
	return nil
}

func (n *JsonParseNode) GetOnAfterAssignFieldValues() absser.ParsableAction {
	return n.onAfterAssignFieldValues
}

func (n *JsonParseNode) SetOnAfterAssignFieldValues(action absser.ParsableAction) error {
	n.onAfterAssignFieldValues = action
	return nil
}

// rawToUntypedNodeable converts a raw primitive value (stored without a JsonParseNode wrapper)
// to an absser.UntypedNodeable. Used when constructing untyped node trees from optimised
// (allocation-reduced) parse trees.
func rawToUntypedNodeable(v interface{}) absser.UntypedNodeable {
	switch rv := v.(type) {
	case *bool:
		return absser.NewUntypedBoolean(*rv)
	case *string:
		return absser.NewUntypedString(*rv)
	case *float32:
		return absser.NewUntypedFloat(*rv)
	case *float64:
		return absser.NewUntypedDouble(*rv)
	case *int32:
		return absser.NewUntypedInteger(*rv)
	case *int64:
		return absser.NewUntypedLong(*rv)
	case *int8:
		return absser.NewUntypedInteger(int32(*rv))
	case *byte:
		return absser.NewUntypedInteger(int32(*rv))
	default:
		return absser.NewUntypedNode(v)
	}
}
