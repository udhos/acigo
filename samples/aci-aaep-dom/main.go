package main

import (
	"fmt"
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 2 {
		log.Fatalf("usage: %s vmware-add|vmware-del|l3-add|l3-del|l2-add|l2-del|list args", os.Args[0])
	}

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing
	list, errList := a.AttachableAccessEntityProfileList()
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}
	for _, t := range list {
		name := t["name"]

		log.Printf("aaep=%s", name)

		aaep, isStr := name.(string)
		if !isStr {
			log.Printf("  aaep name is not a string: %v", name)
		}

		domains, errDom := a.AttachableAccessEntityProfileDomainList(aaep)
		if errDom != nil {
			log.Printf("  could not list domains: %v", errDom)
			continue
		}

		for _, d := range domains {
			dom := d["tDn"]
			log.Printf("  domain=%s", dom)
		}
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "vmware-add":
		if len(args) < 2 {
			log.Fatalf("usage: %s vmware-add aaep dom-vmware", os.Args[0])
		}
		aaep := args[0]
		dom := args[1]
		errAdd := a.AttachableAccessEntityProfileDomainVmmVMWareAdd(aaep, dom)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s", aaep, dom)
	case "vmware-del":
		if len(args) < 2 {
			log.Fatalf("usage: %s vmware-del aaep dom-vmware", os.Args[0])
		}
		aaep := args[0]
		dom := args[1]
		errDel := a.AttachableAccessEntityProfileDomainVmmVMWareDel(aaep, dom)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", aaep, dom)
	case "l3-add":
		if len(args) < 2 {
			log.Fatalf("usage: %s l3-add aaep dom-vmware", os.Args[0])
		}
		aaep := args[0]
		dom := args[1]
		errAdd := a.AttachableAccessEntityProfileDomainL3Add(aaep, dom)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s", aaep, dom)
	case "l3-del":
		if len(args) < 2 {
			log.Fatalf("usage: %s l3-del aaep dom-vmware", os.Args[0])
		}
		aaep := args[0]
		dom := args[1]
		errDel := a.AttachableAccessEntityProfileDomainL3Del(aaep, dom)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", aaep, dom)
	case "l2-add":
		if len(args) < 2 {
			log.Fatalf("usage: %s l2-add aaep dom-vmware", os.Args[0])
		}
		aaep := args[0]
		dom := args[1]
		errAdd := a.AttachableAccessEntityProfileDomainL2Add(aaep, dom)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s", aaep, dom)
	case "l2-del":
		if len(args) < 2 {
			log.Fatalf("usage: %s l2-del aaep dom-vmware", os.Args[0])
		}
		aaep := args[0]
		dom := args[1]
		errDel := a.AttachableAccessEntityProfileDomainL2Del(aaep, dom)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", aaep, dom)
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
