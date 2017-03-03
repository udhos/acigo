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
	list, errList := a.VmmDomainVMWareCredentialsList(domain)
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}
	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		user := t["usr"]
		descr := t["descr"]
		log.Printf("FOUND VMM Domain VMWare Credentials credentials=%s dn=%s user=%s descr=%s", name, dn, user, descr)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 5 {
			log.Fatalf("usage: %s add domain credentials descr user password", os.Args[0])
		}
		domain := args[0]
		credentials := args[1]
		descr := args[2]
		user := args[3]
		password := args[4]
		errAdd := a.VmmDomainVMWareCredentialsAdd(domain, credentials, descr, user, password)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s %s", domain, credentials, descr, user, password)
	case "del":
		if len(args) < 2 {
			log.Fatalf("usage: %s del domain credentials", os.Args[0])
		}
		domain := args[0]
		credentials := args[1]
		errDel := a.VmmDomainVMWareCredentialsDel(domain, credentials)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", domain, credentials)
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
