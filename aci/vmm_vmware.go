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
