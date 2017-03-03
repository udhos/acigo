package aci

import (
	"bytes"
	"fmt"
)

func rnVmmDomain(domain string) string {
	return "dom-" + domain
}

// VmmDomainVMWareAdd creates a VMWare VMM Domain.
func (c *Client) VmmDomainVMWareAdd(domain string) error {

	me := "VmmDomainVMWareAdd"

	rn := rnVmmDomain(domain)

	api := "/api/node/mo/uni/vmmp-VMware/" + rn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vmmDomP":{"attributes":{"dn":"uni/vmmp-VMware/%s","name":"%s","rn":"%s","status":"created"},"children":[{"vmmVSwitchPolicyCont":{"attributes":{"dn":"uni/vmmp-VMware/%s/vswitchpolcont","status":"created,modified"}}}]}}`,
		rn, domain, rn, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VmmDomainVMWareDel deletes a VMWare VMM Domain.
func (c *Client) VmmDomainVMWareDel(domain string) error {

	me := "VmmDomainVMWareDel"

	rn := rnVmmDomain(domain)

	api := "/api/node/mo/uni/vmmp-VMware/" + rn + ".json"

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, err := c.delete(url)
	if err != nil {
		return fmt.Errorf("%s: %v", me, err)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VmmDomainVMWareList retrieves the list of VMWare VMM Domains.
func (c *Client) VmmDomainVMWareList() ([]map[string]interface{}, error) {

	me := "VmmDomainVMWareList"

	key := "compDom"

	api := "/api/node/mo/comp/prov-VMware.json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// VmmDomainVMWareVlanPoolSet sets the VLAN pool for the VMWare VMM domain.
func (c *Client) VmmDomainVMWareVlanPoolSet(domain, vlanpool, vlanpoolMode string) error {

	me := "VmmDomainVMWareVlanPoolSet"

	rnD := rnVmmDomain(domain)
	rn := nameVP(vlanpool, vlanpoolMode)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + ".json"

	url := c.getURL(api)

	/*
		j := fmt.Sprintf(`{"vmmDomP":{"attributes":{"dn":"uni/vmmp-VMware/%s","name":"%s","rn":"%s","status":"modified"},"children":[{"infraRsVlanNs":{"attributes":{"tDn":"uni/infra/%s","status":"modified"},"children":[]}},{"vmmVSwitchPolicyCont":{"attributes":{"dn":"uni/vmmp-VMware/%s/vswitchpolcont","status":"created,modified"}}}]}}`,
			rnD, domain, rnD, rn, rnD)
	*/
	j := fmt.Sprintf(`{"vmmDomP":{"attributes":{"dn":"uni/vmmp-VMware/%s","name":"%s","rn":"%s","status":"modified"},"children":[{"infraRsVlanNs":{"attributes":{"tDn":"uni/infra/%s","status":"modified"}}}]}}`,
		rnD, domain, rnD, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VmmDomainVMWareVlanPoolGet retrieves the VLAN pool for the VMWare VMM domain.
func (c *Client) VmmDomainVMWareVlanPoolGet(domain string) (string, string, error) {

	me := "VmmDomainVMWareVlanPoolGet"

	key := "infraRsVlanNs"

	rnD := rnVmmDomain(domain)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return "", "", fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	attrs, errAttr := jsonImdataAttributes(c, body, key, me)
	if errAttr != nil {
		return "", "", errAttr
	}

	if len(attrs) != 1 {
		return "", "", fmt.Errorf("%s: bad attr count=%d", me, len(attrs))
	}

	attr := attrs[0]

	d, found := attr["tDn"]
	if !found {
		return "", "", fmt.Errorf("%s: vlanpool attribute not found", me)
	}

	value, isStr := d.(string)
	if !isStr {
		return "", "", fmt.Errorf("%s: vlanpool attribute is not a string", me)
	}

	if value == "" {
		return "", "", fmt.Errorf("%s: empty vlanpool", me)
	}

	tail := extractTail(value)
	pool, mode := vlanpoolSplit(tail)

	return pool, mode, nil
}
