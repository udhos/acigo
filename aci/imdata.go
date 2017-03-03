package aci

import (
	"encoding/json"
	"fmt"
)

// imdataExtractError looks for error possibly returned in imdata structure
func imdataExtractError(reply interface{}) error {

	imdata, imdataError := mapGet(reply, "imdata")
	if imdataError != nil {
		return fmt.Errorf("imdata error: %v", imdataError)
	}

	list, isList := imdata.([]interface{})
	if !isList {
		return fmt.Errorf("imdata does not hold a list")
	}

	if len(list) == 0 {
		return nil // ok
	}

	first := list[0]

	e, errErr := mapGet(first, "error")
	if errErr != nil {
		return nil // ok
	}

	attr := mapSimple(e, "attributes")
	code := mapString(attr, "code")
	text := mapString(attr, "text")

	return fmt.Errorf("imdata error found: code=%s text=%s", code, text)
}

func parseJSONError(body []byte) error {

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return errJSON
	}

	return imdataExtractError(reply) // return error, if any
}

func jsonImdataAttributes(c hasDebugf, body []byte, key, label string) ([]map[string]interface{}, error) {

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return nil, errJSON
	}

	return imdataAttributes(c, reply, key, label)
}

func imdataAttributes(c hasDebugf, reply interface{}, key, label string) ([]map[string]interface{}, error) {

	imdata, errImdata := mapGet(reply, "imdata")
	if errImdata != nil {
		return nil, fmt.Errorf("%s: missing imdata: %v", label, errImdata)
	}

	list, isList := imdata.([]interface{})
	if !isList {
		return nil, fmt.Errorf("%s: imdata does not hold a list", label)
	}

	return extractKeyAttributes(c, list, key, label), nil
}

func extractKeyAttributes(c hasDebugf, list []interface{}, key, label string) []map[string]interface{} {

	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		item, errItem := mapGet(i, key)
		if errItem != nil {
			c.debugf("%s: not a %s: %v", label, key, i)
			continue
		}
		attr, errAttr := mapGet(item, "attributes")
		if errAttr != nil {
			c.debugf("%s: missing attributes: %v", label, item)
			continue
		}
		m, isMap := attr.(map[string]interface{})
		if !isMap {
			c.debugf("%s: not a map: %v", label, attr)
			continue
		}
		result = append(result, m)
	}

	return result
}

type hasDebugf interface {
	debugf(fmt string, v ...interface{})
}
