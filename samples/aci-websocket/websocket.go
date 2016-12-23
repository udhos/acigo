package main

import (
	"fmt"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	a, errNew := aci.New(aci.ClientOptions{Debug: debug})
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

	fmt.Printf("login: ok\n")

	if errSock := a.WebsocketOpen(); errSock != nil {
		fmt.Printf("websocket error: %v\n", errSock)
		return
	}

	fmt.Printf("open notification websocket: ok\n")

	subscriptionId, errSub := a.TenantSubscribe()
	if errSub != nil {
		fmt.Printf("tenant subscribe error: %v\n", errSub)
		return
	}

	fmt.Printf("subscribe to tenant notifications: ok\n")

	errAdd := a.TenantAdd("tenant-example", "")
	if errAdd != nil {
		fmt.Printf("tenant add error: %v\n", errAdd)
		return
	}

	fmt.Printf("create tenant: ok\n")

	errDel := a.TenantDel("tenant-example")
	if errDel != nil {
		fmt.Printf("tenant del error: %v\n", errDel)
		return
	}

	fmt.Printf("delete tenant: ok\n")

	var msg interface{}
	if errRead := a.WebsocketReadJson(&msg); errRead != nil {
		fmt.Printf("ERROR: websocket read: %v\n", errRead)
		return
	}

	fmt.Printf("SUCCESS: websocket message: %v\n", msg)

	errSubRefresh := a.TenantSubscriptionRefresh(subscriptionId)
	if errSubRefresh != nil {
		fmt.Printf("tenant subscription refresh error: %v", errSubRefresh)
		return
	}

	fmt.Printf("refresh subscription: ok\n")

	errLogout := a.Logout()
	if errLogout != nil {
		fmt.Printf("logout error: %v\n", errLogout)
		return
	}

	fmt.Printf("logout: ok\n")
}
