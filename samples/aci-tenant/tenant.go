package main

import (
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatalf("usage: %s add|del tenant [description]", os.Args[0])
	}

	cmd := os.Args[1]
	name := os.Args[2]
	var descr string
	if len(os.Args) > 3 {
		descr = os.Args[3]
	}

	a := login()
	defer logout(a)

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
	default:
		log.Printf("unknown command: %s", cmd)
		return
	}
}

func login() *aci.Client {
	a, errNew := aci.New(aci.ClientOptions{Debug: true})
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
