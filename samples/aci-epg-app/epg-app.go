package main

import (
	"fmt"
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 4 {
		log.Fatalf("usage: %s add|del|list args", os.Args[0])
	}

	tenant := os.Args[2]
	ap := os.Args[3]

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	list, errList := a.ApplicationEPGList(tenant, ap)
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}

	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		descr := t["descr"]

		log.Printf("FOUND application EPG: name=%s dn=%s descr=%s", name, dn, descr)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 4 {
			log.Fatalf("usage: %s add tenant application-profile bridge-domain epg [descr]", os.Args[0])
		}
		tenant := args[0]
		ap := args[1]
		bd := args[2]
		epg := args[3]
		var descr string
		if len(args) > 4 {
			descr = args[4]
		}
		errAdd := a.ApplicationEPGAdd(tenant, ap, bd, epg, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s %s", tenant, ap, bd, epg, descr)
	case "del":
		if len(args) < 3 {
			log.Fatalf("usage: %s del tenant application-profile epg", os.Args[0])
		}
		tenant := args[0]
		ap := args[1]
		epg := args[2]
		errDel := a.ApplicationEPGDel(tenant, ap, epg)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s %s", tenant, ap, epg)
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
