package main

import (
	"fmt"
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 5 {
		log.Fatalf("usage: %s add|del|list|prov-add|prov-del|cons-add|cons-del args", os.Args[0])
	}

	tenant := os.Args[2]
	ap := os.Args[3]
	epg := os.Args[4]

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing
	{
		list, errList := a.EPGContractProvidedList(tenant, ap, epg)
		if errList != nil {
			log.Printf("could not list: %v", errList)
			return
		}
		for _, t := range list {
			contract := t["tnVzBrCPName"]
			dn := t["tDn"]
			log.Printf("FOUND provided contract=%s dn=%s", contract, dn)
		}
	}
	{
		list, errList := a.EPGContractConsumedList(tenant, ap, epg)
		if errList != nil {
			log.Printf("could not list: %v", errList)
			return
		}
		for _, t := range list {
			contract := t["tnVzBrCPName"]
			dn := t["tDn"]
			log.Printf("FOUND consumed contract=%s dn=%s", contract, dn)
		}
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "prov-add":
		if len(args) < 4 {
			log.Fatalf("usage: %s prov-add tenant application-profile epg contract", os.Args[0])
		}
		tenant := args[0]
		ap := args[1]
		epg := args[2]
		contract := args[3]
		errAdd := a.EPGContractProvidedAdd(tenant, ap, epg, contract)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s", tenant, ap, epg, contract)
	case "prov-del":
		if len(args) < 3 {
			log.Fatalf("usage: %s prov-del tenant application-profile epg contract", os.Args[0])
		}
		tenant := args[0]
		ap := args[1]
		epg := args[2]
		contract := args[3]
		errDel := a.EPGContractProvidedDel(tenant, ap, epg, contract)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s %s %s", tenant, ap, epg, contract)
	case "cons-add":
		if len(args) < 4 {
			log.Fatalf("usage: %s cons-add tenant application-profile epg contract", os.Args[0])
		}
		tenant := args[0]
		ap := args[1]
		epg := args[2]
		contract := args[3]
		errAdd := a.EPGContractConsumedAdd(tenant, ap, epg, contract)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s", tenant, ap, epg, contract)
	case "cons-del":
		if len(args) < 3 {
			log.Fatalf("usage: %s cons-del tenant application-profile epg contract", os.Args[0])
		}
		tenant := args[0]
		ap := args[1]
		epg := args[2]
		contract := args[3]
		errDel := a.EPGContractConsumedDel(tenant, ap, epg, contract)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s %s %s", tenant, ap, epg, contract)
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
