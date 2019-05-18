package config

import (
	"fmt"

	"github.com/gandrille/go-commons/strpair"
)

type entry struct {
	kind       string
	key        string
	parameters []strpair.StrPair
}

func (e entry) hasParameter(key string) bool {
	for _, p := range e.parameters {
		if p.Str1() == key {
			return true
		}
	}
	return false
}

func (e entry) getParameterValue(key string) string {
	for _, p := range e.parameters {
		if p.Str1() == key {
			return p.Str2()
		}
	}
	return ""
}

func (e entry) getParameterSlice(key string) []string {
	var list []string
	for _, p := range e.parameters {
		if p.Str1() == key && p.Str2() != "" {
			list = append(list, p.Str2())
		}
	}
	return list
}

func (e entry) ParameterKeys() []string {
	var keys []string
	for _, p := range e.parameters {
		keys = append(keys, p.Str1())
	}
	return keys
}

func (e entry) checkParams(mandatory, optional []string) error {
	for _, key := range mandatory {
		if !e.hasParameter(key) {
			return fmt.Errorf("%s entry is missing parameter %s", e.kind, key)
		}
	}

	for _, key := range e.ParameterKeys() {
		if !contains(mandatory, key) && !contains(optional, key) {
			return fmt.Errorf("%s entry is having extra parameter %s", e.kind, key)
		}
	}

	return nil
}

func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
