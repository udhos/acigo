package aci

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ClientOptions is used to specify options for the Client.
type ClientOptions struct {
	Hosts []string // List of apic hostnames. If unspecified, env var APIC_HOSTS is used.
	User  string   // Username. If unspecified, env var APIC_USER is used.
	Pass  string   // Password. If unspecified, env var APIC_PASS is used.
	Debug bool     // Debug enables verbose debugging messages to console.
}

// Client is an instance for interacting with ACI using API calls.
type Client struct {
	Opt                 ClientOptions   // Options for the APIC client
	host                int             // Index for current host
	cli                 *http.Client    // Client context for HTTP
	loginToken          string          // Save APIC login token
	loginRefreshTimeout time.Duration   // Save APIC refresh period
	loginRefreshLast    time.Time       // Save APIC last refresh
	socket              *websocket.Conn // APIC websocket for receiving notifications
}

// Environment variables used as default parameters.
const (
	ApicHosts = "APIC_HOSTS" // Env var. List of apic hostnames. Example: "1.1.1.1" or "1.1.1.1,2.2.2.2,3.3.3.3" or "apic1,4.4.4.4"
	ApicUser  = "APIC_USER"  // Env var. Username. Example: "joe"
	ApicPass  = "APIC_PASS"  // Env var. Password. Example: "joesecret"
)

const (
	contentTypeJSON = "application/json" // ACI API ignores Content-Type, but we set it rightly anyway
)

// New creates a new Client instance for interacting with ACI using API calls.
func New(o ClientOptions) (*Client, error) {
	if len(o.Hosts) < 1 {
		hosts := os.Getenv(ApicHosts)
		if hosts == "" {
			return nil, fmt.Errorf("missing apic hosts: %s=%s", ApicHosts, o.Hosts)
		}
		o.Hosts = strings.Split(hosts, ",")
		if len(o.Hosts) < 1 {
			return nil, fmt.Errorf("missing apic hosts: %s=%s", ApicHosts, o.Hosts)
		}
		for _, h := range o.Hosts {
			if strings.TrimSpace(h) == "" {
				return nil, fmt.Errorf("blank apic hostname '%s' in %s=%s", h, ApicHosts, o.Hosts)
			}
		}
	}

	if o.User == "" {
		o.User = os.Getenv(ApicUser)
		if o.User == "" {
			return nil, fmt.Errorf("missing apic user: %s=%s", ApicUser, o.User)
		}
	}

	if o.Pass == "" {
		o.Pass = os.Getenv(ApicPass)
		if o.Pass == "" {
			return nil, fmt.Errorf("missing apic pass: %s=%s", ApicPass, o.Pass)
		}
	}

	c := &Client{Opt: o}

	c.newHTTPClient()

	c.debugf("new client: hosts=%s user=%s pass=%s", c.Opt.Hosts, c.Opt.User, c.Opt.Pass)

	return c, nil
}

func (c *Client) debugf(fmt string, v ...interface{}) {
	if c.Opt.Debug {
		c.logf("debug "+fmt, v...)
	}
}

func (c *Client) logf(fmt string, v ...interface{}) {
	log.Printf("aci client: "+fmt, v...)
}

func (c *Client) jsonAaaUser() string {
	return fmt.Sprintf(`{"aaaUser": {"attributes": {"name": "%s", "pwd": "%s"}}}`, c.Opt.User, c.Opt.Pass)
}

// Logout closes a session to APIC using the API aaaLogout.
func (c *Client) Logout() error {

	api := "/api/aaaLogout.json"

	aaaUser := c.jsonAaaUser()

	url := c.getURL(api)

	c.debugf("logout: url=%s json=%s", url, aaaUser)

	body, errPost := c.post(url, contentTypeJSON, bytes.NewBufferString(aaaUser))
	if errPost != nil {
		return errPost
	}

	c.debugf("logout: reply: %s", string(body))

	return nil
}

// Login opens a new session into APIC using the API aaaLogin.
func (c *Client) Login() error {

	api := "/api/aaaLogin.json"

	aaaUser := c.jsonAaaUser()

	c.debugf("login: api=%s json=%s", api, aaaUser)

	body, errPost := c.postScan(api, contentTypeJSON, bytes.NewBufferString(aaaUser))
	if errPost != nil {
		return errPost
	}

	// Can't get last refresh before .postScan() because
	// .postScan() might hit some broken hosts before connecting.
	refreshLast := time.Now()

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return errJSON
	}

	imdata, imdataError := mapGet(reply, "imdata")
	if imdataError != nil {
		return fmt.Errorf("login: json imdata error: %s", string(body))
	}

	first, firstError := sliceGet(imdata, 0)
	if firstError != nil {
		return fmt.Errorf("login: imdata first error: %s", string(body))
	}

	mm, mmMap := first.(map[string]interface{})
	if !mmMap {
		return fmt.Errorf("login: imdata slice first member not map: %s", string(body))
	}

	for k, v := range mm {
		switch k {
		case "error":
			attr := mapSimple(v, "attributes")
			code := mapString(attr, "code")
			text := mapString(attr, "text")
			return fmt.Errorf("login: error: code=%s text=%s", code, text)
		case "aaaLogin":
			attr := mapSimple(v, "attributes")
			token := mapString(attr, "token")
			refreshTimeout := mapString(attr, "refreshTimeoutSeconds")

			c.saveRefresh(token, refreshTimeout, refreshLast)

			return nil // ok
		}
	}

	return fmt.Errorf("login: could not find aaaLogin response: %s", string(body))
}

// Refresh resets the session timer on APIC using the API aaaRefresh.
// In order to keep the session active, Refresh() must be called at a period lower than the timeout reported by RefreshTimeout().
func (c *Client) Refresh() error {

	api := "/api/aaaRefresh.json"

	url := c.getURL(api)

	refreshLast := time.Now()

	body, errGet := c.get(url)
	if errGet != nil {
		return errGet
	}

	var reply interface{}
	errJSON := json.Unmarshal(body, &reply)
	if errJSON != nil {
		return errJSON
	}

	imdata, imdataError := mapGet(reply, "imdata")
	if imdataError != nil {
		return fmt.Errorf("refresh: json imdata error: %s", string(body))
	}

	first, firstError := sliceGet(imdata, 0)
	if firstError != nil {
		return fmt.Errorf("refresh: imdata first error: %s", string(body))
	}

	mm, mmMap := first.(map[string]interface{})
	if !mmMap {
		return fmt.Errorf("refresh: imdata slice first member not map: %s", string(body))
	}

	for k, v := range mm {
		switch k {
		case "error":
			attr := mapSimple(v, "attributes")
			code := mapString(attr, "code")
			text := mapString(attr, "text")
			return fmt.Errorf("refresh: error: code=%s text=%s", code, text)
		case "aaaLogin":
			attr := mapSimple(v, "attributes")
			token := mapString(attr, "token")
			refreshTimeout := mapString(attr, "refreshTimeoutSeconds")

			c.saveRefresh(token, refreshTimeout, refreshLast)

			return nil // ok
		}
	}

	return fmt.Errorf("refresh: could not find aaaLogin response: %s", string(body))
}

func (c *Client) saveRefresh(token, refreshTimeout string, refreshLast time.Time) {
	c.loginToken = token // save token

	timeout, timeoutErr := strconv.Atoi(refreshTimeout)
	if timeoutErr != nil {
		c.logf("saveRefresh: bad refresh timeout '%s': %v", refreshTimeout, timeoutErr)
		timeout = 60 // defaults to 60 seconds
	}
	c.loginRefreshTimeout = time.Duration(timeout) * time.Second // save timeout
	c.loginRefreshLast = refreshLast

	c.debugf("saveRefresh: token=%s timeout=%v deadline=%s", token, c.RefreshTimeout(), c.RefreshDeadline())
}

// RefreshTimeout gets the session timeout reported by last API call to APIC.
// In order to keep the session active, Refresh() must be called at a period lower than the timeout reported by RefreshTimeout().
func (c *Client) RefreshTimeout() time.Duration {
	return c.loginRefreshTimeout
}

// RefreshDeadline gets the deadline for session timeout.
// In order to keep the session active, Refresh() must be called before that deadline.
func (c *Client) RefreshDeadline() time.Time {
	return c.loginRefreshLast.Add(c.loginRefreshTimeout)
}

func tlsConfig() *tls.Config {
	return &tls.Config{
		CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true,
		MaxVersion:               tls.VersionTLS12,
		MinVersion:               tls.VersionTLS11,
	}
}

func (c *Client) newHTTPClient() {
	tr := &http.Transport{
		TLSClientConfig:    tlsConfig(),
		DisableCompression: true,
		DisableKeepAlives:  true,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c.cli = &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}
}

// getURL builds HTTPS URL for API access.
func (c *Client) getURL(api string) string {
	return makeURL("https", c.Opt.Hosts[c.host], api)
}

// getURLws builds websocket URL for notifications.
func (c *Client) getURLws(api string) string {
	return makeURL("wss", c.Opt.Hosts[c.host], api)
}

// url builds URL from protocol, host, path.
func makeURL(proto, host, path string) string {
	return proto + "://" + host + path
}

// postScan scans multiple APIC hosts.
func (c *Client) postScan(api string, contentType string, r io.Reader) ([]byte, error) {
	var last error

	if isURL(api) {
		return nil, fmt.Errorf("bad api=%s", api)
	}

	// Reset to first APIC host, if all APIC hosts have been exhausted.
	if c.host == len(c.Opt.Hosts) {
		c.host = 0
	}

	for ; c.host < len(c.Opt.Hosts); c.host++ {

		url := c.getURL(api)

		body, errPost := c.post(url, contentType, r)
		if errPost != nil {
			c.debugf("postScan: error: apic: %s: %v", url, errPost)
			last = errPost
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("no more apic hosts to try - last: %v", last)
}

func (c *Client) showCookies(urlStr string) {
	if c.cli.Jar == nil {
		c.debugf("no cookies to send")
		return
	}

	u, errURL := url.Parse(urlStr)
	if errURL != nil {
		c.debugf("showCookies: %s: %v", urlStr, errURL)
		return
	}

	cookies := c.cli.Jar.Cookies(u)
	if len(cookies) < 1 {
		c.debugf("no cookies to send url=%s", u)
		return
	}

	for _, ck := range cookies {
		c.debugf("cookie to send: %s", ck.Name)
	}
}

func (c *Client) learnCookies(resp *http.Response) error {
	cookies := resp.Cookies()
	for _, ck := range cookies {
		c.debugf("learnCookies: seen: url=%s cookie=%s", resp.Request.URL, ck.Name)
		if ck.Name == "APIC-cookie" {
			if c.cli.Jar == nil {
				var errNew error
				c.cli.Jar, errNew = cookiejar.New(nil) // new jar
				if errNew != nil {
					return errNew
				}
			}
			c.cli.Jar.SetCookies(resp.Request.URL, []*http.Cookie{ck}) // add single cookie to jar
			c.debugf("learnCookies: learnt: url=%s cookie=%s value=%s", resp.Request.URL, ck.Name, ck.Value)
			break
		}
	}
	return nil
}

func (c *Client) post(url string, contentType string, r io.Reader) ([]byte, error) {
	c.debugf("post: apic endpoint: %s", url)

	if !isURL(url) {
		return nil, fmt.Errorf("bad URL=%s", url)
	}

	c.showCookies(url)

	resp, errPost := c.cli.Post(url, contentType, r)
	if errPost != nil {
		return nil, errPost
	}

	if errLearn := c.learnCookies(resp); errLearn != nil {
		return nil, errLearn
	}

	body, errBody := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, errBody
}

func (c *Client) get(url string) ([]byte, error) {
	c.debugf("get: apic endpoint: %s", url)

	if !isURL(url) {
		return nil, fmt.Errorf("bad URL=%s", url)
	}

	c.showCookies(url)

	resp, errPost := c.cli.Get(url)
	if errPost != nil {
		return nil, errPost
	}

	if errLearn := c.learnCookies(resp); errLearn != nil {
		return nil, errLearn
	}

	body, errBody := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, errBody
}

func isURL(url string) bool {
	return strings.HasPrefix(url, "https://")
}

func (c *Client) delete(url string) ([]byte, error) {
	c.debugf("delete: apic endpoint: %s", url)

	if !isURL(url) {
		return nil, fmt.Errorf("bad URL=%s", url)
	}

	c.showCookies(url)

	req, errNew := http.NewRequest("DELETE", url, nil)
	if errNew != nil {
		return nil, errNew
	}

	resp, errDel := c.cli.Do(req)
	if errDel != nil {
		return nil, errDel
	}

	if errLearn := c.learnCookies(resp); errLearn != nil {
		return nil, errLearn
	}

	body, errBody := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, errBody
}
