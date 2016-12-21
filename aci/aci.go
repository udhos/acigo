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

func (c *Client) Login() error {

	loginApi := "/api/aaaLogin.json"

	loginJson := fmt.Sprintf(`{"aaaUser": {"attributes": {"name": "%s", "pwd": "%s"}}}`, c.Opt.User, c.Opt.Pass)

	if c.Opt.Debug {
		log.Printf("login: api=%s json=%s", loginApi, loginJson)
	}

	body, errPost := c.post(loginApi, "application/json", bytes.NewBufferString(loginJson))
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
			if c.Opt.Debug {
				log.Printf("login: ok: refresh=%s token=%s", refresh, token)
			}
			timeout, timeoutErr := strconv.Atoi(refresh)
			if timeoutErr != nil {
				log.Printf("login: bad refresh timeout '%s': %v", refresh, timeoutErr)
				timeout = 60 // defaults to 60 seconds
			}
			c.loginToken = token // save
			c.loginRefreshTimeout = time.Duration(timeout) * time.Second
			return nil // ok
		}
	}

	return fmt.Errorf("login: could not find aaaLogin response: %s", string(body))
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

func (c *Client) post(api string, contentType string, r io.Reader) ([]byte, error) {
	var last error

	for ; c.host < len(c.Opt.Hosts); c.host++ {

		url := "https://" + c.Opt.Hosts[c.host] + api

		if c.Opt.Debug {
			log.Printf("trying: apic: %s", url)
		}

		resp, errPost := c.cli.Post(url, contentType, r)
		if errPost != nil {
			if c.Opt.Debug {
				log.Printf("post form error: apic: %s: %v", url, errPost)
			}
			last = errPost
			continue
		}

		body, errBody := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if errBody != nil {
			if c.Opt.Debug {
				log.Printf("body error: apic: %s: %v", url, errBody)
			}
			last = errBody
			continue
		}

		if c.Opt.Debug {
			log.Printf("apic: %s - body=[%v]", url, string(body))
		}

		return body, nil
	}

	return nil, fmt.Errorf("no more apic hosts to try - last: %v", last)
}
