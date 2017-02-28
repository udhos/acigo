package aci

import (
	"bytes"
	"fmt"
)

func rnSubject(subject string) string {
	return "subj-" + subject
}

func dnSubject(tenant, contract, subject string) string {
	return dnContract(tenant, contract) + "/" + rnSubject(subject)
}

// ContractSubjectAdd creates a new subject.
func (c *Client) ContractSubjectAdd(tenant, contract, subject, reverseFilterPorts, descr string) error {

	me := "ContractSubjectAdd"

	rn := rnSubject(subject)
	dn := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	var attrReverse string
	if reverseFilterPorts != "" {
		attrReverse = fmt.Sprintf(`,"revFltPorts":"%s"`, reverseFilterPorts)
	}

	j := fmt.Sprintf(`{"vzSubj":{"attributes":{"dn":"uni/%s","name":"%s","descr":"%s"%s,"rn":"%s","status":"created"}}}`,
		dn, subject, descr, attrReverse, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ContractSubjectDel deletes an existing subject.
func (c *Client) ContractSubjectDel(tenant, contract, subject string) error {

	me := "ContractSubjectDel"

	dnC := dnContract(tenant, contract)
	dn := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dnC + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vzBrCP":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"vzSubj":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		dnC, dn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ContractSubjectList retrieves the list of subjects.
func (c *Client) ContractSubjectList(tenant, contract string) ([]map[string]interface{}, error) {

	me := "ContractSubjectList"

	key := "vzSubj"

	dn := dnContract(tenant, contract)

	api := "/api/node/mo/uni/" + dn + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
