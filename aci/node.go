package aci

import (
	"encoding/json"
	"fmt"
)

// NodeList retrieves the list of top level system elements (APICs, spines, leaves).
func (c *Client) NodeList() ([]map[string]interface{}, error) {

	api := "/api/class/topSystem.json"

	url := c.getURL(api)

	c.debugf("NodeList: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, errGet
	}

	c.debugf("NodeList: reply: %s", string(body))

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return nil, errJSON
	}

	imdata, errImdata := mapGet(reply, "imdata")
	if errImdata != nil {
		return nil, fmt.Errorf("missing imdata: %v", errImdata)
	}

	list, isList := imdata.([]interface{})
	if !isList {
		return nil, fmt.Errorf("imdata does not hold a list: %s", string(body))
	}

	result := make([]map[string]interface{}, 0, len(list))

	for _, n := range list {
		node, errNode := mapGet(n, "topSystem")
		if errNode != nil {
			c.debugf("NodeList: not a topSystem: %v", n)
			continue
		}
		attr, errAttr := mapGet(node, "attributes")
		if errAttr != nil {
			c.debugf("NodeList: missing attributes: %v", node)
			continue
		}
		m, isMap := attr.(map[string]interface{})
		if !isMap {
			c.debugf("NodeList: not a map: %v", attr)
			continue
		}
		result = append(result, m)
	}

	return result, nil
}
