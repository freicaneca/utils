package ruleprocessor

const (
	quoteMark = '\''
	spaceChar = ' '

	opEqual              = "=="
	opNotEqual           = "!="
	opHas                = "has"
	opNotHas             = "!has"
	opGreaterThan        = ">"
	opLowerThan          = "<"
	opGreaterThanOrEqual = ">="
	opLowerThanOrEqual   = "<="
)

type testFunction func(a, b interface{}) bool
type int64TestFunction func(a, b int64) bool

func asInt64(v any) (int64, bool) {
	switch t := v.(type) {
	case int:
		//return int64(v.(int)), true
		return int64(t), true
	case int8:
		return int64(t), true
	case int16:
		return int64(t), true
	case int32:
		return int64(t), true
	case int64:
		return t, true
	case uint:
		return int64(t), true
	case uint8:
		return int64(t), true
	case uint16:
		return int64(t), true
	case uint32:
		return int64(t), true
	case uint64:
		// because of the range of this type (0 to 2^64 -1), the value may
		// get truncated
		return int64(t), true
	}
	return -1, false
}
