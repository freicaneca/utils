package ruleprocessor

import "utils/logging"

type configField struct {
	name  string
	value any
}

func (cf *configField) AsString(log *logging.Logger) string {
	l := log.New()
	v, ok := cf.value.(string)
	if !ok {
		l.Fatal("[%v] is not a string: %#v", cf.name, cf.value)
	}
	return v
}

func (cf *configField) AsInterfaceArray(
	log *logging.Logger,
) []any {
	l := log.New()
	v, ok := cf.value.([]any)
	if !ok {
		l.Fatal("[%v] is not an array of interfaces: %#v", cf.name, cf.value)
	}
	return v
}

func (cf *configField) AsInterfaceMap(
	log *logging.Logger,
) map[string]any {
	l := log.New()
	v, ok := cf.value.(map[string]any)
	if !ok {
		l.Fatal("[%v] is not map of interfaces: %#v", cf.name, cf.value)
	}
	return v
}

func (cf *configField) AsBool(log *logging.Logger) bool {
	l := log.New()
	v, ok := cf.value.(bool)
	if !ok {
		l.Fatal("[%v] is not a boolean: %#v", cf.name, cf.value)
	}
	return v
}

func getField(
	log *logging.Logger,
	name string,
	fields map[string]any,
) *configField {
	l := log.New()
	v, ok := fields[name]
	if !ok {
		l.Fatal("field not found: [%v]", name)
	}
	return &configField{name: name, value: v}
}
