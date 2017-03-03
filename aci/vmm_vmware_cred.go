package aci

import (
	"bytes"
	"fmt"
)

func rnCredentials(credentials string) string {
	return "usracc-" + credentials
}

// VmmDomainVMWareCredentialsAdd creates vCenter Credentials for VMWare VMM Domain.
func (c *Client) VmmDomainVMWareCredentialsAdd(domain, credentials, descr, user, password string) error {

	me := "VmmDomainVMWareCredentialsAdd"

	rnD := rnVmmDomain(domain)
	rn := rnCredentials(credentials)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + "/" + rn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vmmUsrAccP":{"attributes":{"dn":"uni/vmmp-VMware/%s/%s","name":"%s","descr":"%s","usr":"%s","pwd":"%s","rn":"%s","status":"created"}}}`,
		rnD, rn, credentials, descr, user, password, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VmmDomainVMWareCredentialsDel deletes vCenter Credentials from VMWare VMM Domain.
func (c *Client) VmmDomainVMWareCredentialsDel(domain, credentials string) error {

	me := "VmmDomainVMWareCredentialsDel"

	rnD := rnVmmDomain(domain)
	rn := rnCredentials(credentials)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"vmmDomP":{"attributes":{"dn":"uni/vmmp-VMware/%s","status":"modified"},"children":[{"vmmUsrAccP":{"attributes":{"dn":"uni/vmmp-VMware/%s/%s","status":"deleted"}}}]}}`,
		rnD, rnD, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// VmmDomainVMWareCredentialsList retrieves the list of vCenter Credentials in VMWare VMM Domain.
func (c *Client) VmmDomainVMWareCredentialsList(domain string) ([]map[string]interface{}, error) {

	me := "VmmDomainVMWareCredentialsList"

	key := "vmmUsrAccP"

	rnD := rnVmmDomain(domain)

	api := "/api/node/mo/uni/vmmp-VMware/" + rnD + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
