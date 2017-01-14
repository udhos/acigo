package aci

import (
	"bytes"
	"fmt"
)

func domPhysRN(name string) string {
	return "phys-" + name
}

func apiDomain(rn string) string {
	return "/api/node/mo/uni/" + rn + ".json"
}

// PhysicalDomainAdd creates a new physical domain.
func (c *Client) PhysicalDomainAdd(name, vlanpoolName, vlanpoolMode string) error {

	pool := nameVP(vlanpoolName, vlanpoolMode)

	rn := domPhysRN(name)

	api := apiDomain(rn)

	j := fmt.Sprintf(`{"physDomP":{"attributes":{"dn":"uni/%s","name":"%s","rn":"%s","status":"created"},"children":[{"infraRsVlanNs":{"attributes":{"tDn":"uni/infra/%s","status":"created"}}}]}}`,
		rn, name, rn, pool)

	url := c.getURL(api)

	c.debugf("PhysicalDomainAdd: url=%s json=%s", url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return errPost
	}

	c.debugf("PhysicalDomainAdd: reply: %s", string(body))

	return parseJSONError(body)
}

// PhysicalDomainDel deletes an existing physical domain.
func (c *Client) PhysicalDomainDel(name string) error {

	rn := domPhysRN(name)

	api := apiDomain(rn)

	url := c.getURL(api)

	c.debugf("PhysicalDomainAdd: url=%s", url)

	body, errDel := c.delete(url)
	if errDel != nil {
		return errDel
	}

	c.debugf("PhysicalDomainDel: reply: %s", string(body))

	return parseJSONError(body)
}

// PhysicalDomainList retrieves the list of physical domains.
func (c *Client) PhysicalDomainList() ([]map[string]interface{}, error) {

	key := "physDomP"

	api := "/api/node/mo/uni.json?query-target=subtree&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("PhysicalDomainList: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, errGet
	}

	c.debugf("PhysicalDomainList: reply: %s", string(body))

	return jsonImdataAttributes(c, body, key, "PhysicalDomainList")
}

// PhysicalDomainVlanPoolGet retrieves the VLAN pool for the physical domain.
func (c *Client) PhysicalDomainVlanPoolGet(name string) (string, error) {

	key := "infraRsVlanNs"

	rn := domPhysRN(name)

	api := "/api/node/mo/uni/" + rn + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("PhysicalDomainVlanPoolGet: url=%s", url)

	body, errDel := c.get(url)
	if errDel != nil {
		return "", errDel
	}

	c.debugf("PhysicalDomainVlanPoolGet: reply: %s", string(body))

	attrs, errAttr := jsonImdataAttributes(c, body, key, "PhysicalDomainVlanPoolGet")
	if errAttr != nil {
		return "", errAttr
	}

	if len(attrs) < 1 {
		return "", fmt.Errorf("empty list of vlanpool")
	}

	attr := attrs[0]
	pool := attr["tDn"]

	poolName, isStr := pool.(string)
	if !isStr {
		return "", fmt.Errorf("vlanpool is not a string")
	}

	return poolName, nil
}
