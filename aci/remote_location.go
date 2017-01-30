package aci

import (
	"bytes"
	"fmt"
)

func rnPath(location string) string {
	return "path-" + location
}

// RemoteLocationAdd creates a new remote location.
func (c *Client) RemoteLocationAdd(location, host, protocol, remotePort, remotePath, username, password, descr string) error {

	me := "RemoteLocationAdd"

	rn := rnPath(location)

	api := "/api/node/mo/uni/fabric/" + rn + ".json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fileRemotePath":{"attributes":{"dn":"uni/fabric/%s","remotePort":"%s","name":"%s","descr":"%s","host":"%s","protocol":"%s","remotePath":"%s","userName":"%s","userPasswd":"%s","rn":"%s","status":"created"}}}`,
		rn, remotePort, location, descr, host, protocol, remotePath, username, password, rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// RemoteLocationDel deletes an existing remote location.
func (c *Client) RemoteLocationDel(location string) error {

	me := "RemoteLocationDel"

	rn := rnPath(location)

	api := "/api/node/mo/uni/fabric.json"

	url := c.getURL(api)

	j := fmt.Sprintf(`{"fabricInst":{"attributes":{"dn":"uni/fabric","status":"modified"},"children":[{"fileRemotePath":{"attributes":{"dn":"uni/fabric/%s","status":"deleted"}}}]}}`,
		rn)

	c.debugf("%s: url=%s json=%s", me, url, j)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(j))
	if errPost != nil {
		return fmt.Errorf("%s: %v", me, errPost)
	}

	c.debugf("%s: reply: %s", me, string(body))

	return parseJSONError(body)
}

// RemoteLocationList retrieves the list of remote locations.
func (c *Client) RemoteLocationList() ([]map[string]interface{}, error) {

	me := "RemoteLocationList"

	key := "fileRemotePath"

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
