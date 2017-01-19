package main

import (
	"fmt"
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 3 {
		log.Fatalf("usage: %s add|del|list|vrf-set|vrf-get args", os.Args[0])
	}

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	tenant := os.Args[2]

	list, errList := a.BridgeDomainList(tenant)
	if errList != nil {
		log.Printf("could not list bridge domains: %v", errList)
		return
	}

	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		mac := t["mac"]
		descr := t["descr"]
		log.Printf("found bridge domain: name=%s dn=%s mac=%s descr=%s\n", name, dn, mac, descr)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 2 {
			log.Fatalf("usage: %s add tenant bridge-domain [descr]", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		var descr string
		if len(args) > 2 {
			descr = args[2]
		}
		errAdd := a.BridgeDomainAdd(tenant, bd, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s", tenant, bd)
	case "del":
		if len(args) < 2 {
			log.Fatalf("usage: %s del tenant bridge-domain", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		errDel := a.BridgeDomainDel(tenant, bd)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", tenant, bd)
	case "list":
	case "vrf-set":
		if len(args) < 3 {
			log.Fatalf("usage: %s vrf-set tenant bridge-domain vrf", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		vrf := args[2]
		errDel := a.BridgeDomainVrfSet(tenant, bd, vrf)
		if errDel != nil {
			log.Printf("FAILURE: vrf-set error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: vrf-set: tenant=%s bd=%s vrf=%s", tenant, bd, vrf)
	case "vrf-get":
		if len(args) < 2 {
			log.Fatalf("usage: %s vrf-get tenant bridge-domain", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		vrf, errDel := a.BridgeDomainVrfGet(tenant, bd)
		if errDel != nil {
			log.Printf("FAILURE: vrf-set error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: vrf-get: tenant=%s bd=%s: => vrf=%s", tenant, bd, vrf)
	default:
		log.Printf("unknown command: %s", cmd)
	}
}

func login(debug bool) (*aci.Client, error) {

	a, errNew := aci.New(aci.ClientOptions{Debug: debug})
	if errNew != nil {
		return nil, fmt.Errorf("login new client error: %v", errNew)
	}

	errLogin := a.Login()
	if errLogin != nil {
		return nil, fmt.Errorf("login error: %v", errLogin)
	}

	return a, nil
}

func logout(a *aci.Client) {
	errLogout := a.Logout()
	if errLogout != nil {
		log.Printf("logout error: %v", errLogout)
		return
	}
}
