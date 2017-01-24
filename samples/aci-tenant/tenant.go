package main

import (
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 2 {
		log.Fatalf("usage: %s add|del|list tenant [description]", os.Args[0])
	}

	cmd := os.Args[1]
	var name, descr string
	if len(os.Args) > 2 {
		name = os.Args[2]
	}
	if len(os.Args) > 3 {
		descr = os.Args[3]
	}

	a := login(debug)
	defer logout(a)

	// add/del tenants

	execute(a, cmd, name, descr)

	// display existing tenants

	tenants, errList := a.TenantList()
	if errList != nil {
		log.Printf("could not list tenants: %v", errList)
		return
	}

	for _, t := range tenants {
		name := t["name"]
		dn := t["dn"]
		descr := t["descr"]
		log.Printf("FOUND tenant: name=%s dn=%s descr=%s\n", name, dn, descr)
	}
}

func execute(a *aci.Client, cmd, name, descr string) {
	switch cmd {
	case "add":
		errAdd := a.TenantAdd(name, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add tenant: %s", name)
	case "del":
		errDel := a.TenantDel(name)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del tenant: %s", name)
	case "list":
	default:
		log.Printf("unknown command: %s", cmd)
	}
}

func login(debug bool) *aci.Client {

	a, errNew := aci.New(aci.ClientOptions{Debug: debug})
	if errNew != nil {
		log.Printf("login new client error: %v", errNew)
		os.Exit(1)
	}

	errLogin := a.Login()
	if errLogin != nil {
		log.Printf("login error: %v", errLogin)
		os.Exit(1)
	}

	return a
}

func logout(a *aci.Client) {
	errLogout := a.Logout()
	if errLogout != nil {
		log.Printf("logout error: %v", errLogout)
		return
	}

	log.Printf("logout: done")
}
