package aci

import (
	"bytes"
	"fmt"
)

func nameVR(from, to string) string {
	return "from-[vlan-" + from + "]-to-[vlan-" + to + "]"
}

// VlanRangeAdd creates a new VLAN range for a VLAN pool.
func (c *Client) VlanRangeAdd(vlanpoolName, vlanpoolMode, from, to string) error {

	pool := nameVP(vlanpoolName, vlanpoolMode)

	rang := nameVR(from, to)

	api := "/api/node/mo/uni/infra/" + pool + "/" + rang + ".json"

	j := fmt.Sprintf(`{"fvnsEncapBlk":{"attributes":{"dn":"uni/infra/%s/%s","from":"vlan-%s","to":"vlan-%s","rn":"%s","status":"created"}}}`,
		pool, rang, from, to, rang)

	url := c.getURL(api)

	c.debugf("VlanRangeAdd: url=%s json=%s", url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return errPost
	}

	c.debugf("VlanRangeAdd: reply: %s", string(body))

	return parseJSONError(body)
}

// VlanRangeDel deletes an existing VLAN range from a VLAN pool.
func (c *Client) VlanRangeDel(vlanpoolName, vlanpoolMode, from, to string) error {

	pool := nameVP(vlanpoolName, vlanpoolMode)

	rang := nameVR(from, to)

	api := "/api/node/mo/uni/infra/" + pool + ".json"

	j := fmt.Sprintf(`{"fvnsVlanInstP":{"attributes":{"dn":"uni/infra/%s","status":"modified"},"children":[{"fvnsEncapBlk":{"attributes":{"dn":"uni/infra/%s/%s","status":"deleted"}}}]}}`,
		pool, pool, rang)

	url := c.getURL(api)

	c.debugf("VlanRangeAdd: url=%s json=%s", url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return errPost
	}

	c.debugf("VlanRangeDel: reply: %s", string(body))

	return parseJSONError(body)
}

// VlanRangeList retrieves the list of VLAN ranges from a VLAN pool.
func (c *Client) VlanRangeList(vlanpoolName, vlanpoolMode string) ([]map[string]interface{}, error) {

	pool := nameVP(vlanpoolName, vlanpoolMode)

	key := "fvnsEncapBlk"

	api := "/api/node/mo/uni/infra/" + pool + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("VlanRangeList: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, errGet
	}

	c.debugf("VlanRangeList: reply: %s", string(body))

	return jsonImdataAttributes(c, body, key, "VlanRangeList")
}
