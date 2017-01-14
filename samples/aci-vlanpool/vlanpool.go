package main

import (
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	var listCmd bool
	if len(os.Args) > 1 {
		listCmd = os.Args[1] == "list"
	}

	var cmd, name, mode, descr string

	if !listCmd {
		if len(os.Args) < 4 {
			log.Fatalf("usage: %s add|del|list vlanpool mode [description]", os.Args[0])
		}

		cmd = os.Args[1]
		name = os.Args[2]
		mode = os.Args[3]
		if len(os.Args) > 4 {
			descr = os.Args[4]
		}

		if mode != "static" && mode != "dynamic" {
			log.Fatalf("bad mode=%s: expecting 'static' or 'dynamic'", mode)
		}
	}

	a := login(debug)
	defer logout(a)

	// add/del

	execute(a, cmd, name, mode, descr)

	// display existing

	aps, errList := a.VlanPoolList()
	if errList != nil {
		log.Printf("could not list VLAN pools: %v", errList)
		return
	}

	for _, t := range aps {
		name := t["name"]
		mode := t["allocMode"]
		dn := t["dn"]
		descr := t["descr"]
		log.Printf("found VLAN pool: name=%s mode=%s dn=%s descr=%s\n", name, mode, dn, descr)
	}
}

func execute(a *aci.Client, cmd, name, mode, descr string) {
	switch cmd {
	case "add":
		errAdd := a.VlanPoolAdd(name, mode, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s", name)
	case "del":
		errDel := a.VlanPoolDel(name, mode)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s", name)
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
