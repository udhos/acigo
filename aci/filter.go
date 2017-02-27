package aci

import (
	"bytes"
	"fmt"
)

func rnFilter(filter string) string {
	return "flt-" + filter
}

func dnFilter(tenant, filter string) string {
	return rnTenant(tenant) + "/" + rnFilter(filter)
}

// FilterAdd creates a new filter.
func (c *Client) FilterAdd(tenant, filter, descr string) error {

	me := "FilterAdd"

	rn := rnFilter(filter)
	dn := dnFilter(tenant, filter)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzFilter":{"attributes":{"dn":"uni/%s","name":"%s","descr":"%s","rn":"%s","status":"created,modified"}}}`,
		dn, filter, descr, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// FilterDel deletes an existing filter.
func (c *Client) FilterDel(tenant, filter string) error {

	me := "FilterDel"

	rnT := rnTenant(tenant)
	dn := dnFilter(tenant, filter)

	api := "/api/node/mo/uni/" + rnT + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvTenant":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"vzFilter":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		rnT, dn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// FilterList retrieves the list of filters.
func (c *Client) FilterList(tenant string) ([]map[string]interface{}, error) {

	me := "FilterList"

	key := "vzFilter"

	rnT := rnTenant(tenant)

	api := "/api/node/mo/uni/" + rnT + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
