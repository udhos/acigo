package aci

import (
	"bytes"
	"fmt"
)

func rnNode(serial string) string {
	return "nodep-" + serial
}

func dnNode(serial string) string {
	return "controller/nodeidentpol/" + rnNode(serial)
}

// NodeAdd creates a new fabric membership node.
func (c *Client) NodeAdd(name, ID, serial string) error {

	me := "NodeAdd"

	rn := rnNode(serial)

	dn := dnNode(serial)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fabricNodeIdentP":{"attributes":{"dn":"uni/%s","serial":"%s","nodeId":"%s","name":"%s","rn":"%s","status":"created"},"children":[]}}`,
		dn, serial, ID, name, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// NodeDel deletes an existing fabric membership node.
func (c *Client) NodeDel(serial string) error {

	me := "NodeDel"

	dn := dnNode(serial)

	api := "/api/node/mo/uni/controller/nodeidentpol.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fabricNodeIdentP":{"attributes":{"dn":"uni/%s","status":"deleted"},"children":[]}}`,
		dn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// NodeList retrieves the list of top level system elements (APICs, spines, leaves).
func (c *Client) NodeList() ([]map[string]interface{}, error) {

	key := "topSystem"

	api := "/api/class/" + key + ".json"

	url := c.getURL(api)

	c.debugf("NodeList: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, errGet
	}

	c.debugf("NodeList: reply: %s", string(body))

	return jsonImdataAttributes(c, body, key, "NodeList")
}
