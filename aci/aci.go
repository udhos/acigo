package aci

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	//"net/url"
	"bytes"
	"io"
	"os"
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
	Opt  ClientOptions
	host int
	cli  *http.Client
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
			return nil, fmt.Errorf("missing apic hosts: %s=%s", APIC_HOSTS, o.User)
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

	return c, nil
}

func (c *Client) Login() error {

	loginApi := "/api/aaaLogin.json"

	loginJson := fmt.Sprintf("{'aaaUser': {'attributes': {'name': %s, 'pwd': %s}}}", c.Opt.User, c.Opt.Pass)

	_, err := c.post(loginApi, "application/json", bytes.NewBufferString(loginJson))

	return err
}

func (c *Client) post(api string, contentType string, r io.Reader) ([]byte, error) {
	var last error

	for ; c.host < len(c.Opt.Hosts); c.host++ {

		url := "https://" + c.Opt.Hosts[c.host] + api

		if c.cli == nil {

			if c.Opt.Debug {
				log.Printf("trying: apic: %s", url)
			}

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					//CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
					CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
					//RootCAs:                  pool,
					PreferServerCipherSuites: true,
					InsecureSkipVerify:       true,
					//MaxVersion:               tls.VersionTLS11,
					MinVersion: tls.VersionTLS12,
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

		resp, errPost := c.cli.Post(url, contentType, r)
		if errPost != nil {
			if c.Opt.Debug {
				log.Printf("post form error: apic: %s: %v", url, errPost)
			}
			c.cli = nil
			last = errPost
			continue
		}

		body, errBody := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if errBody != nil {
			if c.Opt.Debug {
				log.Printf("body error: apic: %s: %v", url, errBody)
			}
			c.cli = nil
			last = errBody
			continue
		}

		if c.Opt.Debug {
			log.Printf("apic: %s - body=[%v]", url, body)
		}

		return body, nil
	}

	return nil, fmt.Errorf("no more apic hosts to try - last: %v", last)
}

/*
func forcePort(host, port string) string {
	i := strings.LastIndexByte(host, ':')
	if i < 0 {
		return host + port
	}
	return host
}
*/
