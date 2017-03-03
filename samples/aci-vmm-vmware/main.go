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
		log.Fatalf("usage: %s add|del|list args", os.Args[0])
	}

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing
	list, errList := a.VmmDomainVMWareList()
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}
	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		log.Printf("FOUND VMM Domain VMWare name=%s dn=%s", name, dn)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 1 {
			log.Fatalf("usage: %s add domain", os.Args[0])
		}
		dom := args[0]
		errAdd := a.VmmDomainVMWareAdd(dom)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s", dom)
	case "del":
		if len(args) < 1 {
			log.Fatalf("usage: %s del domain", os.Args[0])
		}
		dom := args[0]
		errDel := a.VmmDomainVMWareDel(dom)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s", dom)
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
