package aci

import (
	"bytes"
	"fmt"
)

func dnFilterEntry(tenant, filter, entry string) string {
	return dnFilter(tenant, filter) + "/" + rnFilterEntry(entry)
}

func rnFilterEntry(entry string) string {
	return "e-" + entry
}

// FilterEntryAdd creates a new filter entry.
func (c *Client) FilterEntryAdd(tenant, filter, entry, etherType, ipProto, srcPortFrom, srcPortTo, dstPortFrom, dstPortTo string) error {

	me := "FilterEntryAdd"

	rn := rnFilterEntry(entry)
	dn := dnFilterEntry(tenant, filter, entry)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzEntry":{"attributes":{"dn":"uni/%s","name":"%s","etherT":"%s","status":"created,modified","prot":"%s","sFromPort":"%s","sToPort":"%s","dFromPort":"%s","dToPort":"%s","rn":"%s"}}}`,
		dn, entry, etherType, ipProto, srcPortFrom, srcPortTo, dstPortFrom, dstPortTo, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// FilterEntryDel deletes an existing filter entry.
func (c *Client) FilterEntryDel(tenant, filter, entry string) error {

	me := "FilterEntryDel"

	dnF := dnFilter(tenant, filter)
	dn := dnFilterEntry(tenant, filter, entry)

	api := "/api/node/mo/uni/" + dnF + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzFilter":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"vzEntry":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		dnF, dn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// FilterEntryList retrieves the list of filter entries.
func (c *Client) FilterEntryList(tenant, filter string) ([]map[string]interface{}, error) {

	me := "FilterEntryList"

	key := "vzEntry"

	dnF := dnFilter(tenant, filter)

	api := "/api/node/mo/uni/" + dnF + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
