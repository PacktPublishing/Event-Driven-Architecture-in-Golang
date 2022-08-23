package di

import (
	"strings"
)

type tracked map[string]int

func (s tracked) add(info depInfo) tracked {
	newList := make(tracked, len(s))

	for k, v := range s {
		newList[k] = v
	}
	newList[info.key] = len(newList)

	return newList
}

func (s tracked) ordered() []string {
	keys := make([]string, len(s))

	for key, i := range s {
		keys[i] = key
	}

	return keys
}

func (s tracked) String() string {
	return strings.Join(s.ordered(), ",")
}
