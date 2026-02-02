package ruleprocessor

import "strings"

var testOperatorsList = map[string]testFunction{
	opEqual:              equalsTest,
	opNotEqual:           withNot(equalsTest),
	opHas:                containsTest,
	opNotHas:             withNot(containsTest),
	opGreaterThan:        withInt64(greaterThanTest),
	opLowerThan:          withInt64(lowerThanTest),
	opGreaterThanOrEqual: withInt64(greaterThanOrEqualTest),
	opLowerThanOrEqual:   withInt64(lowerThanOrEqualTest),
}

func equalsTest(a, b interface{}) bool {
	vA := a
	if intA, ok := asInt64(a); ok {
		vA = intA
	}
	return vA == b
}

func containsTest(a, b any) bool {
	switch a.(type) {
	case string:
		switch b.(type) {
		case string:
			//return strings.Index(a.(string), b.(string)) != -1
			return strings.Contains(a.(string), b.(string))
		}
	}
	return false
}

func greaterThanTest(a, b int64) bool {
	return a > b
}

func greaterThanOrEqualTest(a, b int64) bool {
	return a >= b
}

func lowerThanTest(a, b int64) bool {
	return a < b
}

func lowerThanOrEqualTest(a, b int64) bool {
	return a <= b
}

func withInt64(test int64TestFunction) testFunction {
	return func(a, b interface{}) bool {
		intA, aOK := asInt64(a)
		intB, bOK := asInt64(b)
		if aOK && bOK {
			return test(intA, intB)
		}
		return false
	}
}

func withNot(test testFunction) testFunction {
	return func(a, b any) bool {
		return !test(a, b)
	}
}
