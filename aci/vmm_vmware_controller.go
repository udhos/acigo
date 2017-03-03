package aci

import (
	"bytes"
	"fmt"
)

func rnController(controller string) string {
	return "ctrlr-" + controller
}

// VmmDomainVMWareControllerAdd creates controller for VMWare VMM Domain.
func (c *Client) VmmDomainVMWareControllerAdd(domain, controller, credentials, hostname, datacenter string) error {

	me := "VmmDomainVMWareControllerAdd"

	rnD := rnVmmDomain(domain)
	rnC := rnCredentials(credentials)
	rn := rnController(controller)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + "/" + rn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vmmCtrlrP":{"attributes":{"dn":"uni/vmmp-VMware/%s/%s","name":"%s","hostOrIp":"%s","rootContName":"%s","rn":"%s","status":"created"},"children":[{"vmmRsAcc":{"attributes":{"tDn":"uni/vmmp-VMware/%s/%s","status":"created"}}}]}}`,
		rnD, rn, controller, hostname, datacenter, rn, rnD, rnC)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VmmDomainVMWareControllerDel deletes controller from VMWare VMM Domain.
func (c *Client) VmmDomainVMWareControllerDel(domain, controller string) error {

	me := "VmmDomainVMWareControllerDel"

	rnD := rnVmmDomain(domain)
	rn := rnController(controller)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vmmDomP":{"attributes":{"dn":"uni/vmmp-VMware/%s","status":"modified"},"children":[{"vmmCtrlrP":{"attributes":{"dn":"uni/vmmp-VMware/%s/%s","status":"deleted"}}}]}}`,
		rnD, rnD, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VmmDomainVMWareControllerList retrieves the list of controllers in VMWare VMM Domain.
func (c *Client) VmmDomainVMWareControllerList(domain string) ([]map[string]interface{}, error) {

	me := "VmmDomainVMWareControllerList"

	key := "vmmCtrlrP"

	rnD := rnVmmDomain(domain)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + ".json?query-target=subtree&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}

// VmmDomainVMWareControllerCredentialsGet retrieves controller credentials.
func (c *Client) VmmDomainVMWareControllerCredentialsGet(domain, controller string) (string, error) {

	me := "VmmDomainVMWareControllerCredentialsGet"

	key := "vmmRsAcc"

	rnD := rnVmmDomain(domain)
	rn := rnController(controller)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + "/" + rn + ".json?query-target=children&target-subtree-class=" + key

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
		return "", fmt.Errorf("%s: credentials attribute not found", me)
	}

	cred, isStr := d.(string)
	if !isStr {
		return "", fmt.Errorf("%s: credentials attribute is not a string", me)
	}

	if cred == "" {
		return "", fmt.Errorf("%s: empty credentials", me)
	}

	tail := extractTail(cred)
	suffix := stripPrefix(tail, "usracc-")

	return suffix, nil
}
