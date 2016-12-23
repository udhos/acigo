/*
Go package for interacting with Cisco ACI using API calls.

Example

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

See also

Cisco APIC REST API User Guide: http://www.cisco.com/c/en/us/td/docs/switches/datacenter/aci/apic/sw/1-x/api/rest/b_APIC_RESTful_API_User_Guide.html

APIC Management Information Model Reference: https://developer.cisco.com/media/mim-ref

*/
package aci
