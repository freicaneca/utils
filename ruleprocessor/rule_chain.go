package ruleprocessor

import "utils/logging"

type ruleChain struct {
	label      string
	value      string
	conditions string
	tests      []*ruleCondition
}

func newRuleChain(
	log *logging.Logger,
	label string,
	conditions string,
	forceLowerCase bool,
	value string,
) *ruleChain {
	l := log.New()
	rc := &ruleChain{label: label, conditions: conditions, value: value}
	rc = decodeConditionsString(l, rc, conditions, forceLowerCase)
	return rc
}

func (rc *ruleChain) Label() string {
	return rc.label
}

func (rc *ruleChain) Conditions() string {
	return rc.conditions
}

func (rc *ruleChain) Value() string {
	return rc.value
}

func (rc *ruleChain) add(rt *ruleCondition) *ruleChain {
	rc.tests = append(rc.tests, rt)
	return rc
}

func (rc *ruleChain) IsValid(
	log *logging.Logger,
	fields map[string]any,
) bool {
	l := log.New()
	for _, test := range rc.tests {
		k := test.Key()
		fieldContent, ok := fields[k]
		if !ok {
			l.Warn("RuleChain[%v]: missing expected field: [%v]", rc.label, k)
			if test.stopOnFail {
				return false
			}
			// in case it's not necessary to stop, ignore the missing field
			continue
		}
		if result, stop := test.IsOK(fieldContent); stop {
			return result
		}
	}
	return false
}
