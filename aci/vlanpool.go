package aci

import (
	"bytes"
	"fmt"
)

func jsonVlanPoolAdd(name, mode, descr string) string {

	rn := nameVP(name, mode)

	j := fmt.Sprintf(`{"fvnsVlanInstP":{"attributes":{"dn":"uni/infra/%s","name":"%s","descr":"%s","allocMode":"%s","rn":"%s","status":"created"}}}`, rn, name, descr, mode, rn)

	return j
}

func jsonVlanPoolDel(name, mode string) string {

	rn := nameVP(name, mode)

	j := fmt.Sprintf(`{"infraInfra":{"attributes":{"dn":"uni/infra","status":"modified"},"children":[{"fvnsVlanInstP":{"attributes":{"dn":"uni/infra/%s","status":"deleted"}}}]}}`, rn)

	return j
}

// get vlan pool resource name
func nameVP(name, mode string) string {
	return fmt.Sprintf("vlanns-[%s]-%s", name, mode)
}

// VlanPoolAdd creates a new VLAN pool.
func (c *Client) VlanPoolAdd(name, mode, descr string) error {

	rn := nameVP(name, mode)

	api := "/api/node/mo/uni/infra/" + rn + ".json"

	j := jsonVlanPoolAdd(name, mode, descr)

	url := c.getURL(api)

	c.debugf("VlanPoolAdd: url=%s json=%s", url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return errPost
	}

	c.debugf("VlanPoolAdd: reply: %s", string(body))

	return parseJSONError(body)
}

// VlanPoolDel deletes an existing VLAN pool.
func (c *Client) VlanPoolDel(name, mode string) error {

	api := "/api/node/mo/uni/infra.json"

	j := jsonVlanPoolDel(name, mode)

	url := c.getURL(api)

	c.debugf("VlanPoolAdd: url=%s json=%s", url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return errPost
	}

	c.debugf("VlanPoolDel: reply: %s", string(body))

	return parseJSONError(body)
}

// VlanPoolList retrieves the list of VLAN pools.
func (c *Client) VlanPoolList() ([]map[string]interface{}, error) {

	key := "fvnsVlanInstP"

	api := "/api/node/mo/uni/infra.json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("VlanPoolList: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, errGet
	}

	c.debugf("VlanPoolList: reply: %s", string(body))

	return jsonImdataAttributes(c, body, key, "VlanPoolList")
}
