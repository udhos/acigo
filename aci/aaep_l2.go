package aci

import (
	"bytes"
	"fmt"
)

// AttachableAccessEntityProfileDomainL2Add attaches an L2 Domain to the AAEP.
func (c *Client) AttachableAccessEntityProfileDomainL2Add(aep, l2dom string) error {

	me := "AttachableAccessEntityProfileDomainL2Add"

	rnE := rnAEP(aep)
	rn := rnL2Dom(l2dom)

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

// AttachableAccessEntityProfileDomainL2Del detaches an L2 Domain from the AAEP.
func (c *Client) AttachableAccessEntityProfileDomainL2Del(aep, l2dom string) error {

	me := "AttachableAccessEntityProfileDomainL2Del"

	rnE := rnAEP(aep)
	rn := rnL2Dom(l2dom)

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
