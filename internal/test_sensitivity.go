package internal

type TestSensitivity int

const (
	NORMAL_SENSITIVITY TestSensitivity = iota
	PERSONAL_SENSITIVITY
	PRIVATE_SENSITIVITY
	CONFIDENTIAL_SENSITIVITY
)

func (i TestSensitivity) String() string {
	return []string{"normal", "personal", "private", "confidential"}[i]
}
func ParseTestSensitivity(v string) (any, error) {
	result := NORMAL_SENSITIVITY
	switch v {
	case "normal":
		result = NORMAL_SENSITIVITY
	case "personal":
		result = PERSONAL_SENSITIVITY
	case "private":
		result = PRIVATE_SENSITIVITY
	case "confidential":
		result = CONFIDENTIAL_SENSITIVITY
	default:
		return nil, nil
	}
	return &result, nil
}
func SerializeTestSensitivity(values []TestSensitivity) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.String()
	}
	return result
}
func (i TestSensitivity) isMultiValue() bool {
	return false
}
