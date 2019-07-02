package aci

// FaultList retrieves the list of faults in the fabric.
func (c *Client) FaultList() ([]map[string]interface{}, error) {

	key := "faultInst"

	api := "/api/node/class/" + key + ".json"

	url := c.getURL(api)

	c.debugf("FaultList: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, errGet
	}

	c.debugf("FaultList: reply: %s", string(body))

	return jsonImdataAttributes(c, body, key, "FaultList")
}
