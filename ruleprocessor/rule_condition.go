package ruleprocessor

import "strings"

type ruleCondition struct {
	key            string
	value          interface{}
	forceLowerCase bool
	stopOnFail     bool
	doTest         testFunction
}

func newRuleCondition(key string, value interface{}, stopOnFail bool, forceLowerCase bool, test testFunction) *ruleCondition {
	return &ruleCondition{
		key:            key,
		value:          value,
		forceLowerCase: forceLowerCase,
		stopOnFail:     stopOnFail,
		doTest:         test,
	}
}

func (rt *ruleCondition) Key() string {
	return rt.key
}

func (rt *ruleCondition) IsOK(content interface{}) (result bool, doStop bool) {
	if v, ok := content.(string); ok && rt.forceLowerCase {
		content = strings.ToLower(v)
	}
	result = rt.doTest(content, rt.value)
	// AND-based tests must stop the chain on first failure
	// OR-based tests must stop the chain on first success
	//
	// result | stopOnFail | doStop
	// -----------------------------
	// false  | false      | false
	// true   | false      | true
	// false  | true       | true
	// true   | true       | false
	doStop = result != rt.stopOnFail
	return
}
