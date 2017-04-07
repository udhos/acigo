package aci

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// SubjectApplyBothDirections reports whether the subject applies its filters to both directions.
func (c *Client) SubjectApplyBothDirections(tenant, contract, subject string) (bool, error) {

	me := "SubjectApplyBothDirections"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + ".json?query-target=children&target-subtree-class=vzInTerm&target-subtree-class=vzOutTerm"

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return false, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return false, errJSON
	}

	count, errCount := mapGet(reply, "totalCount")
	if errCount != nil {
		return false, fmt.Errorf("%s: totalCount error: %s", me, string(body))
	}

	err := imdataExtractError(reply)
	if err != nil {
		return false, err
	}

	both := count != "1" && count != "2"

	return both, nil
}

// SubjectFilterBothAdd attaches a filter to subject.
// This type of filter is applied to both directions.
func (c *Client) SubjectFilterBothAdd(tenant, contract, subject, filter string) error {

	me := "SubjectFilterBothAdd"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzRsSubjFiltAtt":{"attributes":{"tnVzFilterName":"%s","status":"created,modified"}}}`,
		filter)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// SubjectFilterBothDel detaches a filter from subject.
// This type of filter is applied to both directions.
func (c *Client) SubjectFilterBothDel(tenant, contract, subject, filter string) error {

	me := "SubjectFilterBothDel"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzSubj":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"vzRsSubjFiltAtt":{"attributes":{"dn":"uni/%s/rssubjFiltAtt-%s","status":"deleted"}}}]}}`,
		dnS, dnS, filter)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// SubjectFilterBothList retrieves the list of filters attached to subject.
// These filters are applied to both directions.
func (c *Client) SubjectFilterBothList(tenant, contract, subject string) ([]map[string]interface{}, error) {

	me := "SubjectFilterBothList"

	key := "vzRsSubjFiltAtt"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// SubjectFilterInputAdd attaches an input filter to subject.
func (c *Client) SubjectFilterInputAdd(tenant, contract, subject, filter string) error {

	me := "SubjectFilterInputAdd"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + "/intmnl.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzRsFiltAtt":{"attributes":{"tnVzFilterName":"%s","status":"created"}}}`,
		filter)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// SubjectFilterInputDel detaches an input filter from subject.
func (c *Client) SubjectFilterInputDel(tenant, contract, subject, filter string) error {

	me := "SubjectFilterInputDel"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + "/intmnl.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzInTerm":{"attributes":{"dn":"uni/%s/intmnl","status":"modified"},"children":[{"vzRsFiltAtt":{"attributes":{"dn":"uni/%s/intmnl/rsfiltAtt-%s","status":"deleted"}}}]}}`,
		dnS, dnS, filter)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// SubjectFilterInputList retrieves the list of input filters attached to subject.
func (c *Client) SubjectFilterInputList(tenant, contract, subject string) ([]map[string]interface{}, error) {

	me := "SubjectFilterInputList"

	key := "vzRsFiltAtt"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + "/intmnl.json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// SubjectFilterOutputAdd attaches an output filter to subject.
func (c *Client) SubjectFilterOutputAdd(tenant, contract, subject, filter string) error {

	me := "SubjectFilterOutputAdd"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + "/outtmnl.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzRsFiltAtt":{"attributes":{"tnVzFilterName":"%s","status":"created"}}}`,
		filter)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// SubjectFilterOutputDel detaches an output filter from subject.
func (c *Client) SubjectFilterOutputDel(tenant, contract, subject, filter string) error {

	me := "SubjectFilterOutputDel"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + "/outtmnl.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzInTerm":{"attributes":{"dn":"uni/%s/outtmnl","status":"modified"},"children":[{"vzRsFiltAtt":{"attributes":{"dn":"uni/%s/outtmnl/rsfiltAtt-%s","status":"deleted"}}}]}}`,
		dnS, dnS, filter)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// SubjectFilterOutputList retrieves the list of output filters attached to subject.
func (c *Client) SubjectFilterOutputList(tenant, contract, subject string) ([]map[string]interface{}, error) {

	me := "SubjectFilterOutputList"

	key := "vzRsFiltAtt"

	dnS := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnS + "/outtmnl.json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
