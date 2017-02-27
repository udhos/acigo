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
		log.Fatalf("usage: %s add|del|list|entry-add|entry-del args", os.Args[0])
	}

	tenant := os.Args[2]

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	list, errList := a.FilterList(tenant)
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}

	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		descr := t["descr"]

		log.Printf("FOUND filter: name=%s dn=%s descr=%s", name, dn, descr)

		filter, isStr := name.(string)
		if !isStr {
			log.Printf("  filter=%s not a string", filter)
			continue
		}

		entries, errEntries := a.FilterEntryList(tenant, filter)
		if errEntries != nil {
			log.Printf("  filter=%s could not list entries: %v", filter, errEntries)
			continue
		}

		for _, e := range entries {
			entry := e["name"]
			etherType := e["etherT"]
			ipProto := e["prot"]
			sFromPort := e["sFromPort"]
			sToPort := e["sToPort"]
			dFromPort := e["dFromPort"]
			dToPort := e["dToPort"]
			log.Printf("  filter=%s entry=%s etherType=%s ipProto=%s srcPortFrom=%s srcPortTo=%s dstPortFrom=%s dstPortTo=%s",
				filter, entry, etherType, ipProto, sFromPort, sToPort, dFromPort, dToPort)
		}
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 2 {
			log.Fatalf("usage: %s add tenant filter [descr]", os.Args[0])
		}
		tenant := args[0]
		filter := args[1]
		var descr string
		if len(args) > 2 {
			descr = args[2]
		}
		errAdd := a.FilterAdd(tenant, filter, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s", tenant, filter, descr)
	case "del":
		if len(args) < 2 {
			log.Fatalf("usage: %s del tenant filter", os.Args[0])
		}
		tenant := args[0]
		filter := args[1]
		errDel := a.FilterDel(tenant, filter)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", tenant, filter)
	case "list":
	case "entry-add":
		if len(args) < 9 {
			log.Fatalf("usage: %s add tenant filter entry ether-type ip-proto src-port-from src-port-to dst-port-from dst-port-to", os.Args[0])
		}
		tenant := args[0]
		filter := args[1]
		entry := args[2]
		etherType := args[3]
		ipProto := args[4]
		srcPortFrom := args[5]
		srcPortTo := args[6]
		dstPortFrom := args[7]
		dstPortTo := args[8]
		errAdd := a.FilterEntryAdd(tenant, filter, entry, etherType, ipProto, srcPortFrom, srcPortTo, dstPortFrom, dstPortTo)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s %s %s %s %s %s", tenant, filter, entry, etherType, ipProto, srcPortFrom, srcPortTo, dstPortFrom, dstPortTo)
	case "entry-del":
		if len(args) < 3 {
			log.Fatalf("usage: %s del tenant filter entry", os.Args[0])
		}
		tenant := args[0]
		filter := args[1]
		entry := args[2]
		errDel := a.FilterEntryDel(tenant, filter, entry)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s %s", tenant, filter, entry)
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
