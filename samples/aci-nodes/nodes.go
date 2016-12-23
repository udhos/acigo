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

	nodes, errNodes := a.NodeList()
	if errNodes != nil {
		fmt.Printf("nodes error: %v\n", errNodes)
		return
	}

	for _, n := range nodes {
		dn := n["dn"]
		role := n["role"]
		fmt.Printf("FOUND node: role=%s dn=%s\n", role, dn)
	}

	errLogout := a.Logout()
	if errLogout != nil {
		fmt.Printf("logout error: %v\n", errLogout)
		return
	}

	fmt.Printf("logout: ok\n")
}
