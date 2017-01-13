package aci

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
