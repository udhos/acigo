package aci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
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

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(jsonTenant))
	if errPost != nil {
		return errPost
	}

	c.debugf("tenant add: reply: %s", string(body))

	return parseJSONError(body)
}

func parseJSONError(body []byte) error {

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return errJSON
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

// TenantDel deletes an existing tenant.
func (c *Client) TenantDel(name string) error {

	api := "/api/mo/uni.json"

	jsonTenant := jsonTenantDel(name)

	url := c.getURL(api)

	//url += "?rsp-subtree=modified" // demand response

	c.debugf("tenant del: url=%s json=%s", url, jsonTenant)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(jsonTenant))
	if errPost != nil {
		return errPost
	}

	c.debugf("tenant del: reply: %s", string(body))

	return parseJSONError(body)
}

// TenantSubscribe deletes an existing tenant.
// The subscriptionId is returned.
func (c *Client) TenantSubscribe() (string, error) {

	api := "/api/class/fvTenant.json?subscription=yes"

	url := c.getURL(api)

	c.debugf("tenant subscribe: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return "", errGet
	}

	c.debugf("tenant subscribe: reply: %s", string(body))

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return "", errJSON
	}

	sub, subError := mapGet(reply, "subscriptionId")
	if subError != nil {
		return "", fmt.Errorf("tentant subscribe error %v: %s", subError, string(body))
	}

	subscriptionId, isStr := sub.(string)
	if !isStr {
		return "", fmt.Errorf("subId not a string %v: %s", sub, string(body))
	}

	c.debugf("TenantSubscribe: subscriptionId=%s", subscriptionId)

	return subscriptionId, nil
}

// TenantSubscriptionTimeout gets the subscription timeout.
// In order to keep the subscription active, TenantSubscriptionRefresh() must be called at a period lower than the timeout reported by TenantSubscriptionTimeout().
func (c *Client) TenantSubscriptionTimeout() time.Duration {
	return 60 * time.Second // ACI API docs says this value is fixed
}

// TenantSubscriptionRefresh refreshes a subscription.
// In order to keep the subscription active, TenantSubscriptionRefresh() must be called at a period lower than the timeout reported by TenantSubscriptionTimeout().
func (c *Client) TenantSubscriptionRefresh(subscriptionId string) error {

	api := "/api/subscriptionRefresh.json?id=" + subscriptionId

	url := c.getURL(api)

	c.debugf("TenantSubscriptionRefresh: url=%s", url)

	body, errGet := c.get(url)
	if errGet != nil {
		return errGet
	}

	c.debugf("TenantSubscriptionRefresh: reply: %s", string(body))

	return parseJSONError(body)
}
