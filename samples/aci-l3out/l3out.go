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
		log.Fatalf("usage: %s add|del|list|vrf-set|vrf-get args", os.Args[0])
	}

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	tenant := os.Args[2]

	list, errList := a.L3ExtOutList(tenant)
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}

	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		descr := t["descr"]

		log.Printf("found external routed network: name=%s dn=%s descr=%s", name, dn, descr)

		out, isStr := name.(string)
		if !isStr {
			log.Printf("external routed network name is not string: %s", name)
			continue
		}

		vrf, errVrfGet := a.L3ExtOutVrfGet(tenant, out)
		if errVrfGet == nil {
			log.Printf("  external routed network %s vrf=[%s]", out, vrf)
		}
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 2 {
			log.Fatalf("usage: %s add tenant out [descr]", os.Args[0])
		}
		tenant := args[0]
		out := args[1]
		var descr string
		if len(args) > 2 {
			descr = args[2]
		}
		errAdd := a.L3ExtOutAdd(tenant, out, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s", tenant, out)
	case "del":
		if len(args) < 2 {
			log.Fatalf("usage: %s del tenant out", os.Args[0])
		}
		tenant := args[0]
		out := args[1]
		errDel := a.L3ExtOutDel(tenant, out)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", tenant, out)
	case "list":
	case "vrf-set":
		if len(args) < 3 {
			log.Fatalf("usage: %s vrf-set tenant out vrf", os.Args[0])
		}
		tenant := args[0]
		out := args[1]
		vrf := args[2]
		errAdd := a.L3ExtOutVrfSet(tenant, out, vrf)
		if errAdd != nil {
			log.Printf("FAILURE: vrf-set error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: vrf-set: tenant=%s out=%s vrf=%s", tenant, out, vrf)
	case "vrf-get":
		if len(args) < 2 {
			log.Fatalf("usage: %s vrf-get tenant out", os.Args[0])
		}
		tenant := args[0]
		out := args[1]
		vrf, errGet := a.L3ExtOutVrfGet(tenant, out)
		if errGet != nil {
			log.Printf("FAILURE: vrf-set error: %v", errGet)
			return
		}
		log.Printf("SUCCESS: vrf-get: tenant=%s out=%s: => vrf=%s", tenant, out, vrf)
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
