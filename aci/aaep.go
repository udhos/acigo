package aci

import (
	"bytes"
	"fmt"
)

func rnAEP(aep string) string {
	return "attentp-" + aep
}

// AttachableAccessEntityProfileAdd creates an AAEP.
func (c *Client) AttachableAccessEntityProfileAdd(aep, descr string) error {

	me := "AttachableAccessEntityProfileAdd"

	rn := rnAEP(aep)

	api := "/api/node/mo/uni/infra.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"infraInfra":{"attributes":{"dn":"uni/infra","status":"modified"},"children":[{"infraAttEntityP":{"attributes":{"dn":"uni/infra/%s","name":"%s","descr":"%s","rn":"%s","status":"created"}}}]}}`,
		rn, aep, descr, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// AttachableAccessEntityProfileDel deletes an AAEP.
func (c *Client) AttachableAccessEntityProfileDel(aep string) error {

	me := "AttachableAccessEntityProfileDel"

	rn := rnAEP(aep)

	api := "/api/node/mo/uni/infra.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"infraInfra":{"attributes":{"dn":"uni/infra","status":"modified"},"children":[{"infraAttEntityP":{"attributes":{"dn":"uni/infra/%s","status":"deleted"}}}]}}`,
		rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// AttachableAccessEntityProfileList retrieves the list of AAEPs.
func (c *Client) AttachableAccessEntityProfileList() ([]map[string]interface{}, error) {

	me := "AttachableAccessEntityProfileList"

	key := "infraAttEntityP"

	api := "/api/node/mo/uni/infra.json?query-target=subtree&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
