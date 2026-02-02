package ruleprocessor

import (
	"utils/logging"
)

const (
	defaultValueField   = "default"
	forceLowerCaseField = "force-lowercase"
	rulesField          = "rules"
	labelField          = "label"
	conditionsField     = "conditions"
	valueField          = "value"
)

type Processor struct {
	defaultValue string
	chains       []*ruleChain
}

func NewFrom(
	log *logging.Logger,
	c map[string]any,
) *Processor {
	l := log.New("Processor.NewFrom")
	defaultValue := getField(l, defaultValueField, c).AsString(l)
	forceLowerCase := getField(l, forceLowerCaseField, c).AsBool(l)
	rules := getField(l, rulesField, c).AsInterfaceArray(l)
	p := New(defaultValue)
	for _, rl := range rules {
		r := (&configField{name: "rule", value: rl}).AsInterfaceMap(l)
		label := getField(l, labelField, r).AsString(l)
		conditions := getField(l, conditionsField, r).AsString(l)
		value := getField(l, valueField, r).AsString(l)
		p.AddRule(l, label, conditions, forceLowerCase, value)
	}
	return p
}

func New(defaultValue string) *Processor {
	return &Processor{defaultValue: defaultValue}
}

func (rp *Processor) AddRule(
	log *logging.Logger,
	label string,
	conditions string,
	forceLowerCase bool,
	value string,
) *Processor {
	l := log.New()
	rp.chains = append(rp.chains, newRuleChain(l, label, conditions, forceLowerCase, value))
	return rp
}

func (rp *Processor) Process(
	log *logging.Logger,
	input map[string]any,
) string {
	l := log.New()
	for _, rc := range rp.chains {
		if rc.IsValid(l, input) {
			l.Debug("RuleProcessor: valid rule[%v] => conditions[%q]", rc.Label(), rc.Conditions())
			return rc.Value()
		}
	}
	l.Debug("RuleProcessor: default value[%v]", rp.defaultValue)
	return rp.defaultValue
}
