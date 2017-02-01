package aci

import (
	"bytes"
	"fmt"
)

func rnExportConfig(config string) string {
	return "configexp-" + config
}

// ExportConfigurationRun executes the export configuration now.
func (c *Client) ExportConfigurationRun(config string) error {

	// A policy can be triggered at any time by setting the adminSt to triggered.

	me := "ExportConfigurationRun"

	rn := rnExportConfig(config)

	api := "/api/node/mo/uni/fabric/" + rn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"configExportP":{"attributes":{"dn":"uni/fabric/%s","adminSt":"triggered"}}}`,
		rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ExportConfigurationAdd creates a new export configuration.
func (c *Client) ExportConfigurationAdd(config, scheduler, remoteLocation, descr string) error {

	me := "ExportConfigurationAdd"

	rn := rnExportConfig(config)

	api := "/api/node/mo/uni/fabric/" + rn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"configExportP":{"attributes":{"dn":"uni/fabric/%s","name":"%s","descr":"%s","rn":"%s","status":"created"},"children":[{"configRsExportScheduler":{"attributes":{"tnTrigSchedPName":"%s","status":"created,modified"}}},{"configRsRemotePath":{"attributes":{"tnFileRemotePathName":"%s","status":"created,modified"}}}]}}`,
		rn, config, descr, rn, scheduler, remoteLocation)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ExportConfigurationDel deletes an existing export configuration.
func (c *Client) ExportConfigurationDel(config string) error {

	me := "ExportConfigurationDel"

	rn := rnExportConfig(config)

	api := "/api/node/mo/uni/fabric.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fabricInst":{"attributes":{"dn":"uni/fabric","status":"modified"},"children":[{"configExportP":{"attributes":{"dn":"uni/fabric/%s","status":"deleted"}}}]}}`,
		rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// ExportConfigurationList retrieves the list of export configurations.
func (c *Client) ExportConfigurationList() ([]map[string]interface{}, error) {

	me := "ExportConfigurationList"

	key := "configExportP"

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

// ExportConfigurationSchedulerGet retrieves the scheduler attached to an export configuration.
func (c *Client) ExportConfigurationSchedulerGet(config string) (map[string]interface{}, error) {

	me := "ExportConfigurationSchedulerGet"

	rn := rnExportConfig(config)

	key := "configRsExportScheduler"

	api := "/api/node/mo/uni/fabric/" + rn + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	list, errImdata := jsonImdataAttributes(c, body, key, me)
	if errImdata != nil {
		return nil, fmt.Errorf("%s: %v", me, errImdata)
	}

	if len(list) < 1 {
		return nil, fmt.Errorf("%s: empty list", me)
	}

	return list[0], nil
}

// ExportConfigurationRemoteLocationGet retrieves the remote location attached to an export configuration.
func (c *Client) ExportConfigurationRemoteLocationGet(config string) (map[string]interface{}, error) {

	me := "ExportConfigurationRemoteLocationGet"

	rn := rnExportConfig(config)

	key := "configRsRemotePath"

	api := "/api/node/mo/uni/fabric/" + rn + ".json?query-target=children&target-subtree-class=" + key

	url := c.getURL(api)

	c.debugf("%s: url=%s", me, url)

	body, errGet := c.get(url)
	if errGet != nil {
		return nil, fmt.Errorf("%s: %v", me, errGet)
	}

	c.debugf("%s: reply: %s", me, string(body))

	list, errImdata := jsonImdataAttributes(c, body, key, me)
	if errImdata != nil {
		return nil, fmt.Errorf("%s: %v", me, errImdata)
	}

	if len(list) < 1 {
		return nil, fmt.Errorf("%s: empty list", me)
	}

	return list[0], nil
}
