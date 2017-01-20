package aci

import (
	"bytes"
	"fmt"
)

// BridgeDomainL3ExtOutAdd attaches a new L3 External Outside in a bridge domain.
func (c *Client) BridgeDomainL3ExtOutAdd(tenant, bd, out string) error {

	me := "BridgeDomainL3ExtOutAdd"

	dn := dnBridgeDomain(tenant, bd)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvRsBDToOut":{"attributes":{"tnL3extOutName":"%s","status":"created"}}}`,
		out)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// BridgeDomainL3ExtOutDel detaches an existing L3 External Outside from a bridge domain.
func (c *Client) BridgeDomainL3ExtOutDel(tenant, bd, out string) error {

	me := "BridgeDomainL3ExtOutDel"

	dnBD := dnBridgeDomain(tenant, bd)

	api := "/api/node/mo/uni/" + dnBD + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvBD":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"fvRsBDToOut":{"attributes":{"dn":"uni/%s/rsBDToOut-%s","status":"deleted"}}}]}}`,
		dnBD, dnBD, out)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// BridgeDomainL3ExtOutList retrieves the list of L3 External Outsides attached to a bridge domain.
func (c *Client) BridgeDomainL3ExtOutList(tenant, bd string) ([]map[string]interface{}, error) {

	me := "BridgeDomainL3ExtOutList"

	key := "fvRsBDToOut"

	dnBD := dnBridgeDomain(tenant, bd)

	api := "/api/node/mo/uni/" + dnBD + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
