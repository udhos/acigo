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
)

type ClientOptions struct {
	Hosts []string
	User  string
	Pass  string
	Debug bool
}

type Client struct {
	Opt                 ClientOptions
	host                int
	cli                 *http.Client
	loginToken          string
	loginRefreshTimeout time.Duration
}

const (
	APIC_HOSTS = "APIC_HOSTS"
	APIC_USER  = "APIC_USER"
	APIC_PASS  = "APIC_PASS"
)

func New(o ClientOptions) (*Client, error) {
	if len(o.Hosts) < 1 {
		hosts := os.Getenv(APIC_HOSTS)
		if hosts == "" {
			return nil, fmt.Errorf("missing apic hosts: %s=%s", APIC_HOSTS, o.Hosts)
		}
		o.Hosts = strings.Split(hosts, ",")
		if len(o.Hosts) < 1 {
			return nil, fmt.Errorf("missing apic hosts: %s=%s", APIC_HOSTS, o.Hosts)
		}
		for _, h := range o.Hosts {
			if strings.TrimSpace(h) == "" {
				return nil, fmt.Errorf("blank apic hostname '%s' in %s=%s", h, APIC_HOSTS, o.Hosts)
			}
		}
	}

	if o.User == "" {
		o.User = os.Getenv(APIC_USER)
		if o.User == "" {
			return nil, fmt.Errorf("missing apic user: %s=%s", APIC_USER, o.User)
		}
	}

	if o.Pass == "" {
		o.Pass = os.Getenv(APIC_PASS)
		if o.Pass == "" {
			return nil, fmt.Errorf("missing apic pass: %s=%s", APIC_PASS, o.Pass)
		}
	}

	c := &Client{Opt: o}

	c.newHttpClient()

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

func (c *Client) Login() error {

	loginApi := "/api/aaaLogin.json"

	loginJson := fmt.Sprintf(`{"aaaUser": {"attributes": {"name": "%s", "pwd": "%s"}}}`, c.Opt.User, c.Opt.Pass)

	c.debugf("login: api=%s json=%s", loginApi, loginJson)

	body, errPost := c.postLogin(loginApi, "application/json", bytes.NewBufferString(loginJson))
	if errPost != nil {
		return errPost
	}

	var reply interface{}
	errJson := json.Unmarshal(body, &reply)
	if errJson != nil {
		return errJson
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
			refresh := mapString(attr, "refreshTimeoutSeconds")
			timeout, timeoutErr := strconv.Atoi(refresh)
			if timeoutErr != nil {
				c.logf("login: bad refresh timeout '%s': %v", refresh, timeoutErr)
				timeout = 60 // defaults to 60 seconds
			}
			c.loginToken = token                                         // save
			c.loginRefreshTimeout = time.Duration(timeout) * time.Second // save
			c.debugf("login: ok: refresh=%v token=%s", c.RefreshTimeout(), token)
			return nil // ok
		}
	}

	return fmt.Errorf("login: could not find aaaLogin response: %s", string(body))
}

func (c *Client) Refresh() error {

	refreshApi := "/api/aaaRefresh.json"

	url := c.getURL(refreshApi)

	body, errGet := c.get(url)
	if errGet != nil {
		return errGet
	}

	var reply interface{}
	errJson := json.Unmarshal(body, &reply)
	if errJson != nil {
		return errJson
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
			refresh := mapString(attr, "refreshTimeoutSeconds")
			c.debugf("refresh: ok: refresh=%s token=%s", refresh, token)
			return nil // ok
		}
	}

	return fmt.Errorf("refresh: could not find aaaLogin response: %s", string(body))
}

func (c *Client) RefreshTimeout() time.Duration {
	return c.loginRefreshTimeout
}

func (c *Client) newHttpClient() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			MaxVersion:               tls.VersionTLS11,
			MinVersion:               tls.VersionTLS11,
		},
		DisableCompression: true,
		DisableKeepAlives:  true,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
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

func (c *Client) getURL(api string) string {
	url := "https://" + c.Opt.Hosts[c.host] + api
	return url
}

func (c *Client) postLogin(api string, contentType string, r io.Reader) ([]byte, error) {
	var last error

	for ; c.host < len(c.Opt.Hosts); c.host++ {

		url := c.getURL(api)

		c.debugf("trying: apic: %s", url)

		body, errPost := c.post(url, contentType, r)
		if errPost != nil {
			c.debugf("post error: apic: %s: %v", url, errPost)
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

	u, errUrl := url.Parse(urlStr)
	if errUrl != nil {
		c.debugf("showCookies: %s: %v", urlStr, errUrl)
		return
	}

	cookies := c.cli.Jar.Cookies(u)
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
	c.debugf("post: %s", url)

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
	c.debugf("get: %s", url)

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
