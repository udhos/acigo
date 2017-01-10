package yname

import (
	"fmt"
	"strconv"
	"strings"
)

// SplitFunc is a utility type for a function which takes a path "a/b/c" and splits it to head="a" and tail="b/c".
// The separator is defined within the function.
type SplitFunc func(path string) (head, tail string)

func splitSep(str string, sep byte) (string, string) {
	slash := strings.IndexByte(str, sep)
	if slash < 0 {
		return str, ""
	}
	return str[:slash], str[slash+1:]
}

// GetSep finds a path within a structure composed of nested maps and slices.
// The path separator sep is given explicitly.
func GetSep(doc interface{}, path string, sep byte) (interface{}, error) {
	split := func(path string) (string, string) {
		return splitSep(path, sep)
	}
	return GetSplit(doc, path, split)
}

// GetSplit finds a path within a structure composed of nested maps and slices.
// You must supply the function to split the path.
func GetSplit(doc interface{}, path string, split SplitFunc) (interface{}, error) {
	head, tail := split(path)
	if head == "" {
		return nil, fmt.Errorf("invalid empty path: [%s]", path)
	}

	switch i := doc.(type) {
	case map[interface{}]interface{}:
		child, found := i[head]
		if !found {
			return nil, fmt.Errorf("not found: [%s]", head)
		}
		if tail == "" {
			return child, nil
		}
		return GetSplit(child, tail, split)
	case []interface{}:
		index, errConv := strconv.Atoi(head)
		if errConv != nil {
			return nil, fmt.Errorf("not and index: %s: %v", head, errConv)
		}
		if index < 0 || index >= len(i) {
			return nil, fmt.Errorf("index not found: %d", index)
		}
		child := i[index]
		if tail == "" {
			return child, nil
		}
		return GetSplit(child, tail, split)
	}

	return nil, fmt.Errorf("unsupported type: [%s]: %v", head, doc)
}
