package aci

import (
	"fmt"
)

func mapGet(i interface{}, member string) (interface{}, error) {
	m, isMap := i.(map[string]interface{})
	if !isMap {
		return nil, fmt.Errorf("json mapGet: not a map")
	}
	mem, found := m[member]
	if !found {
		return nil, fmt.Errorf("json mapGet: member [%s] not found", member)
	}
	return mem, nil
}

func sliceGet(i interface{}, member int) (interface{}, error) {
	list, isList := i.([]interface{})
	if !isList {
		return nil, fmt.Errorf("json sliceGet: not a slice")
	}
	if member < 0 || member >= len(list) {
		return nil, fmt.Errorf("json sliceGet: member=%d out-of-bounds", member)
	}
	return list[member], nil
}

func mapSimple(i interface{}, member string) interface{} {
	m, _ := mapGet(i, member)
	return m
}

func mapString(i interface{}, member string) string {
	m := mapSimple(i, member)
	s, isStr := m.(string)
	if isStr {
		return s
	}
	return ""
}
