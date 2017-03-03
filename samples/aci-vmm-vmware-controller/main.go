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
		log.Fatalf("usage: %s add|del|list args", os.Args[0])
	}

	domain := os.Args[2]

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing
	list, errList := a.VmmDomainVMWareControllerList(domain)
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}
	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		hostname := t["hostOrIp"]
		datacenter := t["rootContName"]

		controller, isStr := name.(string)
		if !isStr {
			log.Printf("controller name not a string: %v", name)
		}

		cred, errCred := a.VmmDomainVMWareControllerCredentialsGet(domain, controller)
		if errCred != nil {
			log.Printf("could not get credentials for controller=%s: %v", controller, errCred)
		}

		log.Printf("FOUND VMM Domain VMWare controller controller=%s dn=%s credentials=%s hostname=%s datacenter=%s", name, dn, cred, hostname, datacenter)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 5 {
			log.Fatalf("usage: %s add domain controller credentials hostname datacenter", os.Args[0])
		}
		domain := args[0]
		controller := args[1]
		credentials := args[2]
		hostname := args[3]
		datacenter := args[4]
		errAdd := a.VmmDomainVMWareControllerAdd(domain, controller, credentials, hostname, datacenter)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s %s", domain, controller, credentials, hostname, datacenter)
	case "del":
		if len(args) < 2 {
			log.Fatalf("usage: %s del domain controller", os.Args[0])
		}
		domain := args[0]
		controller := args[1]
		errDel := a.VmmDomainVMWareControllerDel(domain, controller)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", domain, controller)
	case "list":
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
