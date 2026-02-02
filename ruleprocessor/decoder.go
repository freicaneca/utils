package ruleprocessor

import (
	"fmt"
	"strconv"
	"strings"
	"utils/logging"
)

func mapKeyWords(s string) []string {
	quoted := false
	keys := strings.FieldsFunc(s, func(r rune) bool {
		if r == quoteMark {
			quoted = !quoted
		}
		return !quoted && r == spaceChar
	})
	return keys
}

func decodeConditionsString(
	log *logging.Logger,
	rc *ruleChain,
	conditions string,
	forceLowerCase bool,
) *ruleChain {
	l := log.New()
	keyWords := mapKeyWords(conditions)
	index := 0
	next := func() string {
		if index >= len(keyWords) {
			l.Fatal("syntax error: incomplete rule")
		}
		value := keyWords[index]
		index++
		return value
	}
	finished := func() bool {
		return index >= len(keyWords)
	}
	fatal := func(format string, v ...interface{}) {
		l.Fatal("syntax error: " + fmt.Sprintf(format, v...))
	}
	for !finished() {
		// basic components of a rule's condition: key, operator and value
		key := next()
		valueOp := next()
		value := next()
		logicOp := ""
		if !finished() {
			// gets the logic operator between rule conditions
			logicOp = next()
		}
		// gets the type of test that the operator represents
		testType, ok := testOperatorsList[valueOp]
		if !ok {
			fatal("invalid operator: [%v]", valueOp)
		}
		// based on the type of the logic operator, decides whether the test
		// must stop when the test fails or not
		stopOnFail := false
		if logicOp != "" {
			logicOp = strings.ToLower(logicOp)
			switch logicOp {
			case "and":
				stopOnFail = true
			case "or":
			default:
				fatal("invalid logic operator: [%v]", logicOp)
			}
		}
		// decode the test value
		var testValue interface{}
		invalidOp := func(op, value interface{}) {
			fatal("invalid operation [%v] for value type of [%v]", op, value)
		}
		if value[0] == quoteMark {
			// remove quotation marks from string
			value = value[1 : len(value)-1]
			if forceLowerCase {
				value = strings.ToLower(value)
			}
			testValue = value
			// validate allowed operations for this value type
			switch valueOp {
			case opEqual, opNotEqual, opHas, opNotHas:
			default:
				invalidOp(valueOp, value)
			}
		} else if value == "true" {
			testValue = true
			// validate allowed operations for this value type
			switch valueOp {
			case opEqual, opNotEqual:
			default:
				invalidOp(valueOp, value)
			}
		} else if value == "false" {
			testValue = false
			// validate allowed operations for this value type
			switch valueOp {
			case opEqual, opNotEqual:
			default:
				invalidOp(valueOp, value)
			}
		} else {
			// anything that is not one of the previous types is treated
			// as a number
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				fatal("invalid integer value: [%v]: %v", value, err)
			}
			testValue = intValue
			// validate allowed operations for this value type
			switch valueOp {
			case opEqual, opNotEqual, opGreaterThan, opGreaterThanOrEqual, opLowerThan, opLowerThanOrEqual:
			default:
				invalidOp(valueOp, value)
			}
		}
		rc.add(newRuleCondition(
			key,
			testValue,
			stopOnFail,
			forceLowerCase,
			testType,
		))
	}
	return rc
}
