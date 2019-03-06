[![GoDoc](https://godoc.org/github.com/udhos/acigo/aci?status.svg)](http://godoc.org/github.com/udhos/acigo/aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/acigo)](https://goreportcard.com/report/github.com/udhos/acigo)
[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/acigo/blob/master/LICENSE)

# About Acigo

Acigo is a Go package for interacting with Cisco ACI using API calls.

# Install

## Without Modules - Before Go 1.11

    go get github.com/gorilla/websocket
    go get github.com/udhos/acigo
    go install github.com/udhos/acigo/aci

## With Modules - Since Go 1.11

    git clone https://github.com/udhos/acigo
    cd acigo
    go install ./aci

# Usage

Import the package in your program:

    import "github.com/udhos/acigo/aci"

See godoc: http://godoc.org/github.com/udhos/acigo/aci

# Example

    package main
    
    import (
    	"fmt"
    	"github.com/udhos/acigo/aci"
    )
    
    func main() {
    
    	a, errNew := aci.New(aci.ClientOptions{})
    	if errNew != nil {
    		fmt.Printf("login new client error: %v\n", errNew)
    		return
    	}
    
    	// Since credentials have not been specified explicitly under ClientOptions,
    	// Login() will use env vars: APIC_HOSTS=host, APIC_USER=username, APIC_PASS=pwd
    	errLogin := a.Login()
    	if errLogin != nil {
    		fmt.Printf("login error: %v\n", errLogin)
    		return
    	}
    
    	errAdd := a.TenantAdd("tenant-example", "")
    	if errAdd != nil {
    		fmt.Printf("tenant add error: %v\n", errAdd)
    		return
    	}
    
    	errLogout := a.Logout()
    	if errLogout != nil {
    		fmt.Printf("logout error: %v\n", errLogout)
    		return
    	}
    }

# Documentation

Acigo documentation in GoDoc: https://godoc.org/github.com/udhos/acigo/aci

# See Also

[Cisco APIC REST API User Guide](http://www.cisco.com/c/en/us/td/docs/switches/datacenter/aci/apic/sw/1-x/api/rest/b_APIC_RESTful_API_User_Guide.html)

[APIC Management Information Model Reference](https://developer.cisco.com/media/mim-ref)
