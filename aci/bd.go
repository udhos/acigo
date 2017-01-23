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

func rnSubnet(subnet string) string {
	return "subnet-[" + subnet + "]"
}

func dnSubnet(tenant, bd, subnet string) string {
	return dnBridgeDomain(tenant, bd) + "/" + rnSubnet(subnet)
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
	v, found := attr["tnFvCtxName"]
	if !found {
		return "", fmt.Errorf("%s: VRF not found", me)
	}

	vrf, isStr := v.(string)
	if !isStr {
		return "", fmt.Errorf("%s: VRF is not a string", me)
	}

	if vrf == "" {
		return "", fmt.Errorf("%s: empty VRF name", me)
	}

	return vrf, nil
}

// BridgeDomainSubnetAdd creates a new subnet in a bridge domain.
func (c *Client) BridgeDomainSubnetAdd(tenant, bd, subnet, descr string) error {

	me := "BridgeDomainSubnetAdd"

	rnSN := rnSubnet(subnet)

	dnSN := dnSubnet(tenant, bd, subnet)

	api := "/api/node/mo/uni/" + dnSN + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvSubnet":{"attributes":{"dn":"uni/%s","ip":"%s","descr":"%s","rn":"%s","status":"created"}}}`,
		dnSN, subnet, descr, rnSN)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// BridgeDomainSubnetDel deletes an existing subnet from a bridge domain.
func (c *Client) BridgeDomainSubnetDel(tenant, bd, subnet string) error {

	me := "BridgeDomainSubnetDel"

	dnBD := dnBridgeDomain(tenant, bd)

	dnSN := dnSubnet(tenant, bd, subnet)

	api := "/api/node/mo/uni/" + dnBD + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvBD":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"fvSubnet":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		dnBD, dnSN)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// BridgeDomainSubnetList retrieves the list of subnets from a bridge domain.
func (c *Client) BridgeDomainSubnetList(tenant, bd string) ([]map[string]interface{}, error) {

	me := "BridgeDomainSubnetList"

	key := "fvSubnet"

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

// BridgeDomainSubnetGet retrieves specific subnet from a bridge domain.
func (c *Client) BridgeDomainSubnetGet(tenant, bd, subnet string) ([]map[string]interface{}, error) {

	me := "BridgeDomainSubnetGet"

	key := "fvSubnet"

	dnSN := dnSubnet(tenant, bd, subnet)

	api := "/api/node/mo/uni/" + dnSN + ".json"

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// BridgeDomainSubnetScopeSet defines the scope for a bridge domain subnet.
func (c *Client) BridgeDomainSubnetScopeSet(tenant, bd, subnet, scope string) error {

	me := "BridgeDomainSubnetScopeSet"

	dnSN := dnSubnet(tenant, bd, subnet)

	api := "/api/node/mo/uni/" + dnSN + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvSubnet":{"attributes":{"dn":"uni/%s","scope":"%s"}}}`,
		dnSN, scope)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// BridgeDomainSubnetScopeGet retrieves the scope from a bridge domain subnet.
func (c *Client) BridgeDomainSubnetScopeGet(tenant, bd, subnet string) (string, error) {

	me := "BridgeDomainSubnetScopeGet"

	list, errSubnet := c.BridgeDomainSubnetGet(tenant, bd, subnet)
	if errSubnet != nil {
		return "", fmt.Errorf("%s: %v", me, errSubnet)
	}

	if len(list) < 1 {
		return "", fmt.Errorf("%s: empty list of subnets", me)
	}

	attrs := list[0]
	s := attrs["scope"]

	scope, isStr := s.(string)
	if !isStr {
		return "", fmt.Errorf("%s: scope is not a string", me)
	}

	return scope, nil
}
