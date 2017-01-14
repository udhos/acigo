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

	if len(os.Args) < 4 {
		log.Fatalf("usage: %s add|del|list vlanpool mode from to", os.Args[0])
	}

	var cmd, name, mode, from, to string

	cmd = os.Args[1]
	name = os.Args[2]
	mode = os.Args[3]

	if !listCmd {
		if len(os.Args) < 6 {
			log.Fatalf("usage: %s add|del|list vlanpool mode from to", os.Args[0])
		}

		from = os.Args[4]
		to = os.Args[5]

		if mode != "static" && mode != "dynamic" {
			log.Fatalf("bad mode=%s: expecting 'static' or 'dynamic'", mode)
		}
	}

	a := login(debug)
	defer logout(a)

	// add/del

	execute(a, cmd, name, mode, from, to)

	// display existing

	aps, errList := a.VlanRangeList(name, mode)
	if errList != nil {
		log.Printf("could not list VLAN ranges: %v", errList)
		return
	}

	for _, t := range aps {
		dn := t["dn"]
		mode := t["allocMode"]
		from := t["from"]
		to := t["to"]
		log.Printf("found VLAN range: dn=%s allocMode=%s from=%s to=%s\n", dn, mode, from, to)
	}
}

func execute(a *aci.Client, cmd, name, mode, from, to string) {
	switch cmd {
	case "add":
		errAdd := a.VlanRangeAdd(name, mode, from, to)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s", name)
	case "del":
		errDel := a.VlanRangeDel(name, mode, from, to)
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
