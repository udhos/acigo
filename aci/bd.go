package aci

import (
	"bytes"
	"fmt"
)

func rnBridgeDomain(bd string) string {
	return "BD-" + bd
}

func dnBridgeDomain(tenant, bd string) string {
	return rnTenant(tenant) + "/" + rnBridgeDomain(bd)
}

// BridgeDomainAdd creates a new bridge domain in a tenant.
func (c *Client) BridgeDomainAdd(tenant, bd, descr string) error {

	me := "BridgeDomainAdd"

	rn := rnBridgeDomain(bd)

	dn := dnBridgeDomain(tenant, bd)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvBD":{"attributes":{"dn":"uni/%s","name":"%s","descr":"%s","rn":"%s","status":"created"}}}`,
		dn, bd, descr, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// BridgeDomainDel deletes an existing bridge domain from a tenant.
func (c *Client) BridgeDomainDel(tenant, bd string) error {

	me := "BridgeDomainDel"

	rnT := rnTenant(tenant)

	dn := dnBridgeDomain(tenant, bd)

	api := "/api/node/mo/uni/" + rnT + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvTenant":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"fvBD":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		rnT, dn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// BridgeDomainList retrieves the list of bridge domains from a tenant.
func (c *Client) BridgeDomainList(tenant string) ([]map[string]interface{}, error) {

	me := "BridgeDomainList"

	key := "fvBD"

	t := rnTenant(tenant)

	api := "/api/node/mo/uni/" + t + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// BridgeDomainVrfSet defines the VRF for a bridge domain.
func (c *Client) BridgeDomainVrfSet(tenant, bd, vrf string) error {

	me := "BridgeDomainVrfSet"

	dn := dnBridgeDomain(tenant, bd)

	api := "/api/node/mo/uni/" + dn + "/rsctx.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvRsCtx":{"attributes":{"tnFvCtxName":"%s"}}}`,
		vrf)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// BridgeDomainVrfGet retrieves the VRF for a bridge domain.
func (c *Client) BridgeDomainVrfGet(tenant, bd string) (string, error) {

	me := "BridgeDomainVrfGet"

	key := "fvRsCtx"

	dn := dnBridgeDomain(tenant, bd)

	api := "/api/node/mo/uni/" + dn + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return "", fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	attrs, errAttr := jsonImdataAttributes(c, body, key, me)
	if errAttr != nil {
		return "", fmt.Errorf("%s: %v", me, errAttr)
	}

	if len(attrs) < 1 {
		return "", fmt.Errorf("%s: empty list of VRFs", me)
	}

	attr := attrs[0]
	v := attr["tnFvCtxName"]

	vrf, isStr := v.(string)
	if !isStr {
		return "", fmt.Errorf("%s: VRF is not a string", me)
	}

	return vrf, nil
}
