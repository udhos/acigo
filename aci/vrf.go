package aci

import (
	"bytes"
	"fmt"
)

func rnTenant(tenant string) string {
	return "tn-" + tenant
}

func rnVrf(vrf string) string {
	return "ctx-" + vrf
}

func dnVrf(tenant, vrf string) string {
	return rnTenant(tenant) + "/" + rnVrf(vrf)
}

// VrfAdd creates a new VRF in a tenant.
func (c *Client) VrfAdd(tenant, vrf, descr string) error {

	me := "VrfAdd"

	rn := rnVrf(vrf)

	dn := dnVrf(tenant, vrf)

	api := "/api/node/mo/uni/" + dn + ".json"

	j := fmt.Sprintf(`{"fvCtx":{"attributes":{"dn":"uni/%s","name":"%s","descr":"%s","rn":"%s","status":"created"}}}`, dn, vrf, descr, rn)

	url := c.getURL(api)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VrfDel deletes an existing VRF from a tenant.
func (c *Client) VrfDel(tenant, vrf string) error {

	me := "VrfDel"

	rnT := rnTenant(tenant)

	dn := dnVrf(tenant, vrf)

	api := "/api/node/mo/uni/" + rnT + ".json"

	j := fmt.Sprintf(`{"fvTenant":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"fvCtx":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		rnT, dn)

	url := c.getURL(api)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VrfList retrieves the list of VRFs from a tenant.
func (c *Client) VrfList(tenant string) ([]map[string]interface{}, error) {

	me := "VrfList"

	key := "fvCtx"

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

// VrfSetEnforcedMode sets the VRF enforced mode flag
func (c *Client) VrfSetEnforcedMode(tenant, vrf string, enforced bool) error {
	me := "VrfSetEnforced"
	dn := dnVrf(tenant, vrf)
	var enforcedString = "unenforced"
	if enforced {
		enforcedString = "enforced"
	}
	api := "/api/node/mo/uni/" + dn + ".json"
	j := fmt.Sprintf(`{"fvCtx":{"attributes":{"dn":"uni/%s", "pcEnfPref":"%s"}}}`, dn, enforcedString)
	url := c.getURL(api)
	c.debugf("%s: url=%s json=%s", me, url, j)
	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}
