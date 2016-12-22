package aci

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func jsonTenantAdd(name, descr string) string {

	prefix := fmt.Sprintf(`{"fvTenant":{"attributes":{"name":"%s","status":"created"`, name)
	suffix := "}}}"
	var middle string
	if descr != "" {
		middle = fmt.Sprintf(`,"descr":"%s"`, descr)
	}

	return prefix + middle + suffix
}

func jsonTenantDel(name string) string {
	return fmt.Sprintf(`{"fvTenant":{"attributes":{"name":"%s","status":"deleted"}}}`, name)
}

// TenantAdd creates a new tenant.
func (c *Client) TenantAdd(name, descr string) error {

	api := "/api/mo/uni.json"

	jsonTenant := jsonTenantAdd(name, descr)

	url := c.getURL(api)

	c.debugf("tenant add: url=%s json=%s", url, jsonTenant)

	body, errPost := c.post(url, "application/json", bytes.NewBufferString(jsonTenant))
	if errPost != nil {
		return errPost
	}

	c.debugf("tenant add: reply: %s", string(body))

	return parseJsonError(body)
}

func parseJsonError(body []byte) error {

	var reply interface{}
	errJson := json.Unmarshal(body, &reply)
	if errJson != nil {
		return errJson
	}

	imdata, imdataError := mapGet(reply, "imdata")
	if imdataError != nil {
		return fmt.Errorf("imdata error: %s", string(body))
	}

	list, isList := imdata.([]interface{})
	if !isList {
		return fmt.Errorf("imdata does not hold a list: %s", string(body))
	}

	if len(list) == 0 {
		return nil // ok
	}

	first := list[0]

	e, errErr := mapGet(first, "error")
	if errErr != nil {
		return nil // ok
	}

	attr := mapSimple(e, "attributes")
	code := mapString(attr, "code")
	text := mapString(attr, "text")

	return fmt.Errorf("error: code=%s text=%s", code, text)
}

// TenandDel deletes an existing tenant.
func (c *Client) TenantDel(name string) error {

	api := "/api/mo/uni.json"

	jsonTenant := jsonTenantDel(name)

	url := c.getURL(api)

	//url += "?rsp-subtree=modified" // demand response

	c.debugf("tenant del: url=%s json=%s", url, jsonTenant)

	body, errPost := c.post(url, "application/json", bytes.NewBufferString(jsonTenant))
	if errPost != nil {
		return errPost
	}

	c.debugf("tenant del: reply: %s", string(body))

	return parseJsonError(body)
}
