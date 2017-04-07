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
// reverseFilterPorts: "true", "false", "" (empty means default)
// reverseFilterPorts means the ACI will automatically create a rule with opposite ports to allow the return traffic. For example, if you have a filter allowing traffic from src=X to dst=80, by enabling reverseFilterPorts an implicit rule will be added to allow traffic from src=80 to dst=X.
// applyBothDirections means the subject will apply its filters to both directions.
// If applyBothDirections is enabled, use these functions to manage subject filters: SubjectFilterBothAdd(), SubjectFilterBothDel(), SubjectFilterBothList().
// If applyBothDirections is disabled, use these functions to manage subject filters: SubjectFilterInputAdd(), SubjectFilterInputAddDel(), SubjectFilterInputAdd(), SubjectFilterOutputAdd(), SubjectFilterOutputDel(), SubjectFilterOutputList()
func (c *Client) ContractSubjectAdd(tenant, contract, subject, reverseFilterPorts string, applyBothDirections bool, descr string) error {

	me := "ContractSubjectAdd"

	rn := rnSubject(subject)
	dn := dnSubject(tenant, contract, subject)

	api := "/api/node/mo/uni/" + dn + ".json"

	url := c.getURL(api)

	// reverse filter ports?
	var attrReverse string
	if reverseFilterPorts != "" {
		attrReverse = fmt.Sprintf(`,"revFltPorts":"%s"`, reverseFilterPorts)
	}

	// apply both directions?
	var nonBoth string
	if !applyBothDirections {
		nonBoth = fmt.Sprintf(`,"children":[{"vzInTerm":{"attributes":{"dn":"uni/%s/intmnl","status":"created","targetDscp":"64"},"children":[]}},{"vzOutTerm":{"attributes":{"dn":"uni/%s/outtmnl","status":"created","targetDscp":"64"},"children":[]}}]`,
			dn, dn)
	}

	j := fmt.Sprintf(`{"vzSubj":{"attributes":{"dn":"uni/%s","name":"%s","descr":"%s"%s,"rn":"%s","status":"created"}%s}}`,
		dn, subject, descr, attrReverse, rn, nonBoth)

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
