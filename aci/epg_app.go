package aci

import (
	"bytes"
	"fmt"
)

func rnAEPG(epg string) string {
	return "epg-" + epg
}

func dnAEPG(tenant, ap, epg string) string {
	return dnAP(tenant, ap) + "/" + rnAEPG(epg)
}

// ApplicationEPGAdd creates a new application EPG in an application profile and attached to a bridge domain.
func (c *Client) ApplicationEPGAdd(tenant, applicationProfile, bridgeDomain, epg, descr string) error {

	me := "ApplicationEPGAdd"

	rnE := rnAEPG(epg)

	dnE := dnAEPG(tenant, applicationProfile, epg)

	api := "/api/node/mo/uni/" + dnE + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvAEPg":{"attributes":{"dn":"uni/%s","name":"%s","descr":"%s","rn":"%s","status":"created"},"children":[{"fvRsBd":{"attributes":{"tnFvBDName":"%s","status":"created,modified"}}}]}}`, dnE, epg, descr, rnE, bridgeDomain)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ApplicationEPGDel deletes and existing application EPG from an application profile.
func (c *Client) ApplicationEPGDel(tenant, applicationProfile, epg string) error {

	me := "ApplicationEPGDel"

	dnP := dnAP(tenant, applicationProfile)
	dnE := dnAEPG(tenant, applicationProfile, epg)

	api := "/api/node/mo/uni/" + dnP + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fvAp":{"attributes":{"dn":"uni/%s","status":"modified"},"children":[{"fvAEPg":{"attributes":{"dn":"uni/%s","status":"deleted"}}}]}}`,
		dnP, dnE)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ApplicationEPGList retrieves the list of application EPGs in an application profile.
func (c *Client) ApplicationEPGList(tenant, applicationProfile string) ([]map[string]interface{}, error) {

	me := "ApplicationEPGList"

	key := "fvAEPg"

	dnP := dnAP(tenant, applicationProfile)

	api := "/api/node/mo/uni/" + dnP + ".json?query-target=children&target-subtree-class=" + key + `&query-target-filter=eq(fvAEPg.isAttrBasedEPg,"false")`

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return jsonImdataAttributes(c, body, key, me)
}
