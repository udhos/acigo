package aci

import (
	"bytes"
	"fmt"
)

func rnContract(contract string) string {
	return "brc-" + contract
}

func dnContract(tenant, contract string) string {
	return rnTenant(tenant) + "/" + rnContract(contract)
}

// ContractAdd creates a new contract.
func (c *Client) ContractAdd(tenant, contract, scope, descr string) error {

	me := "ContractAdd"

	rn := rnContract(contract)
	dn := dnContract(tenant, contract)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	var attrScope string
	if scope != "" {
		attrScope = fmt.Sprintf(`,"scope":"%s"`, scope)
	}

	j := fmt.Sprintf(`{"vzBrCP":{"attributes":{"dn":"uni/%s","name":"%s"%s,"descr":"%s","rn":"%s","status":"created"}}}`,
		dn, contract, attrScope, descr, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ContractDel deletes an existing contract.
func (c *Client) ContractDel(tenant, contract string) error {

	me := "ContractDel"

	rnT := rnTenant(tenant)
	dn := dnContract(tenant, contract)

	api := "/api/node/mo/uni/" + rnT + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvTenant":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"vzBrCP":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		rnT, dn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ContractList retrieves the list of contracts.
func (c *Client) ContractList(tenant string) ([]map[string]interface{}, error) {

	me := "ContractList"

	key := "vzBrCP"

	rnT := rnTenant(tenant)

	api := "/api/node/mo/uni/" + rnT + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
