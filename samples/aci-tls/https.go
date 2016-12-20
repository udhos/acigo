package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {

	if len(os.Args) != 4 {
		log.Fatalf("usage: %s host user pass", os.Args[0])
	}

	host := os.Args[1]
	user := os.Args[2]
	pass := os.Args[3]

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS11,
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
	c := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}

	url := fmt.Sprintf("https://%s/api/aaaLogin.json", host)

	loginJson := fmt.Sprintf("{'aaaUser': {'attributes': {'name': %s, 'pwd': %s}}}", user, pass)

	resp, errPost := c.Post(url, "application/json", bytes.NewBufferString(loginJson))
	if errPost != nil {
		log.Fatalf("post error: %v", errPost)
	}

	body, errBody := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if errBody != nil {
		log.Fatalf("body error: %v", errBody)
	}

	fmt.Printf("done - body: %s\n", string(body))
}
