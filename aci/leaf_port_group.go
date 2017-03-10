package aci

import (
	"bytes"
	"fmt"
)

func rnLeafPortGroup(group string) string {
	return "accportgrp-" + group
}

// LeafInterfacePolicyGroupAdd creates a policy group for leaf access ports.
func (c *Client) LeafInterfacePolicyGroupAdd(group, descr string) error {

	me := "LeafInterfacePolicyGroupAdd"

	rn := rnLeafPortGroup(group)

	api := "/api/node/mo/uni/infra/funcprof/" + rn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"infraAccPortGrp":{"attributes":{"dn":"uni/infra/funcprof/%s","name":"%s","descr":"%s","rn":"%s","status":"created"}}}`,
		rn, group, descr, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// LeafInterfacePolicyGroupDel deletes a policy group for leaf access ports.
func (c *Client) LeafInterfacePolicyGroupDel(group string) error {

	me := "LeafInterfacePolicyGroupDel"

	rn := rnLeafPortGroup(group)

	api := "/api/node/mo/uni/infra/funcprof.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"infraFuncP":{"attributes":{"dn":"uni/infra/funcprof","status":"modified"},"children":[{"infraAccPortGrp":{"attributes":{"dn":"uni/infra/funcprof/%s","status":"deleted"}}}]}}`,
		rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// LeafInterfacePolicyGroupList retrieves the list of policy groups for leaf access ports.
func (c *Client) LeafInterfacePolicyGroupList() ([]map[string]interface{}, error) {

	me := "LeafInterfacePolicyGroupList"

	key := "infraAccPortGrp"

	api := "/api/node/class/" + key + ".json"

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// LeafInterfacePolicyGroupEntitySet attaches an AAEP to the leaf interface policy group.
func (c *Client) LeafInterfacePolicyGroupEntitySet(group, aep string) error {

	me := "LeafInterfacePolicyGroupEntitySet"

	rnG := rnLeafPortGroup(group)
	rnE := rnAEP(aep)

	api := "/api/node/mo/uni/infra/funcprof/" + rnG + "/rsattEntP.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"infraRsAttEntP":{"attributes":{"tDn":"uni/infra/%s","status":"created,modified"}}}`,
		rnE)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// LeafInterfacePolicyGroupEntityGet gets the AAEP attached to the leaf interface policy group.
func (c *Client) LeafInterfacePolicyGroupEntityGet(group string) (string, error) {

	me := "LeafInterfacePolicyGroupEntityGet"

	key := "infraRsAttEntP"

	rnG := rnLeafPortGroup(group)

	api := "/api/node/mo/uni/infra/funcprof/" + rnG + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return "", fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	attrs, errAttr := jsonImdataAttributes(c, body, key, me)
	if errAttr != nil {
		return "", errAttr
	}

	if len(attrs) != 1 {
		return "", fmt.Errorf("%s: bad attr count=%d", me, len(attrs))
	}

	attr := attrs[0]

	d, found := attr["tDn"]
	if !found {
		return "", fmt.Errorf("%s: AAEP attribute not found", me)
	}

	aep, isStr := d.(string)
	if !isStr {
		return "", fmt.Errorf("%s: AAEP attribute is not a string", me)
	}

	if aep == "" {
		return "", fmt.Errorf("%s: empty AAEP", me)
	}

	tail := extractTail(aep)
	suffix := stripPrefix(tail, "attentp-")

	return suffix, nil
}
