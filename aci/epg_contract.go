package aci

import (
	"bytes"
	"fmt"
)

// EPGContractProvidedAdd attaches contract as provided by EPG.
func (c *Client) EPGContractProvidedAdd(tenant, applicationProfile, epg, contract string) error {

	me := "EPGContractProvidedAdd"

	dnE := dnAEPG(tenant, applicationProfile, epg)

	api := "/api/node/mo/uni/" + dnE + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`payload{"fvRsProv":{"attributes":{"tnVzBrCPName":"%s","status":"created,modified"}}}`,
		contract)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// EPGContractProvidedDel detaches provided contract from EPG.
func (c *Client) EPGContractProvidedDel(tenant, applicationProfile, epg, contract string) error {

	me := "EPGContractProvidedDel"

	dnE := dnAEPG(tenant, applicationProfile, epg)

	api := "/api/node/mo/uni/" + dnE + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvAEPg":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"fvRsProv":{"attributes":{"dn":"uni/%s/rsprov-%s","status":"deleted"}}}]}}`,
		dnE, dnE, contract)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// EPGContractProvidedList retrieves the list of contracts provided by EPG.
func (c *Client) EPGContractProvidedList(tenant, applicationProfile, epg string) ([]map[string]interface{}, error) {

	me := "EPGContractProvidedList"

	key := "fvRsProv"

	dnE := dnAEPG(tenant, applicationProfile, epg)

	api := "/api/node/mo/uni/" + dnE + ".json?query-target=subtree&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// EPGContractConsumedAdd attaches contract as consumed by EPG.
func (c *Client) EPGContractConsumedAdd(tenant, applicationProfile, epg, contract string) error {

	me := "EPGContractConsumedAdd"

	dnE := dnAEPG(tenant, applicationProfile, epg)

	api := "/api/node/mo/uni/" + dnE + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`payload{"fvRsCons":{"attributes":{"tnVzBrCPName":"%s","status":"created,modified"}}}`,
		contract)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// EPGContractConsumedDel detaches consumed contract from EPG.
func (c *Client) EPGContractConsumedDel(tenant, applicationProfile, epg, contract string) error {

	me := "EPGContractConsumedDel"

	dnE := dnAEPG(tenant, applicationProfile, epg)

	api := "/api/node/mo/uni/" + dnE + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvAEPg":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"fvRsCons":{"attributes":{"dn":"uni/%s/rscons-%s","status":"deleted"}}}]}}`,
		dnE, dnE, contract)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// EPGContractConsumedList retrieves the list of contracts consumed by EPG.
func (c *Client) EPGContractConsumedList(tenant, applicationProfile, epg string) ([]map[string]interface{}, error) {

	me := "EPGContractConsumedList"

	key := "fvRsCons"

	dnE := dnAEPG(tenant, applicationProfile, epg)

	api := "/api/node/mo/uni/" + dnE + ".json?query-target=subtree&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
