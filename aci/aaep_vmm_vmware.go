package aci

import (
	"bytes"
	"fmt"
)

// AttachableAccessEntityProfileDomainVmmVMWareAdd attaches a VMM VMWare Domain to the AAEP.
func (c *Client) AttachableAccessEntityProfileDomainVmmVMWareAdd(aep, domainVMWare string) error {

	me := "AttachableAccessEntityProfileDomainVmmVMWareAdd"

	rnE := rnAEP(aep)
	rn := rnVmmDomainVMWare(domainVMWare)

	api := "/api/node/mo/uni/infra/" + rnE + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"infraRsDomP":{"attributes":{"tDn":"uni/%s","status":"created"}}}`,
		rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// AttachableAccessEntityProfileDomainVmmVMWareDel detaches a VMM VMWare Domain from the AAEP.
func (c *Client) AttachableAccessEntityProfileDomainVmmVMWareDel(aep, domainVMWare string) error {

	me := "AttachableAccessEntityProfileDomainVmmVMWareDel"

	rnE := rnAEP(aep)
	rn := rnVmmDomainVMWare(domainVMWare)

	api := "/api/node/mo/uni/infra/" + rnE + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"infraAttEntityP":{"attributes":{"dn":"uni/infra/%s","status":"modified"},"children":[{"infraRsDomP":{"attributes":{"dn":"uni/infra/%s/rsdomP-[uni/%s]","status":"deleted"}}}]}}`,
		rnE, rnE, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// AttachableAccessEntityProfileDomainList retrieves the list of domains attached to the AAEP.
func (c *Client) AttachableAccessEntityProfileDomainList(aep string) ([]map[string]interface{}, error) {

	me := "AttachableAccessEntityProfileDomainList"

	key := "infraRsDomP"

	rnV := rnAEP(aep)

	api := "/api/node/mo/uni/infra/" + rnV + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
