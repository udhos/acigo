package aci

import (
	"bytes"
	"fmt"
)

// AttachableAccessEntityProfileDomainL3Add attaches an L3 Domain to the AAEP.
func (c *Client) AttachableAccessEntityProfileDomainL3Add(aep, l3dom string) error {

	me := "AttachableAccessEntityProfileDomainL3Add"

	rnE := rnAEP(aep)
	rn := rnL3Dom(l3dom)

	api := "/api/node/mo/uni/infra/" + rnE + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"infraRsDomP":{"attributes":{"tDn":"uni/%s","status":"created"}}}`, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// AttachableAccessEntityProfileDomainL3Del detaches an L3 Domain from the AAEP.
func (c *Client) AttachableAccessEntityProfileDomainL3Del(aep, l3dom string) error {

	me := "AttachableAccessEntityProfileDomainL3Del"

	rnE := rnAEP(aep)
	rn := rnL3Dom(l3dom)

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
