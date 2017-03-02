package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 4 {
		log.Fatalf("usage: %s add|del|list args", os.Args[0])
	}

	tenant := os.Args[2]
	contract := os.Args[3]

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	list, errList := a.ContractSubjectList(tenant, contract)
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}

	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		reverseFilterPorts := t["revFltPorts"]
		descr := t["descr"]

		subject, isStr := name.(string)
		if !isStr {
			log.Printf("subject name not a string: %v", name)
			continue
		}

		var applyBoth string
		both, errBoth := a.SubjectApplyBothDirections(tenant, contract, subject)
		if errBoth != nil {
			log.Printf("subject=%s could not query both directions: %v", subject, errBoth)
			applyBoth = "error"
		} else {
			applyBoth = fmt.Sprintf("%v", both)
		}

		log.Printf("FOUND subject: name=%s dn=%s reverseFilterPorts=%s applyBothDirections=%s descr=%s", name, dn, reverseFilterPorts, applyBoth, descr)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 5 {
			log.Fatalf("usage: %s add tenant contract subject reverse-filter-ports apply-both-directions [descr]", os.Args[0])
		}
		tenant := args[0]
		contract := args[1]
		subject := args[2]
		reverse := args[3]
		both := args[4]
		var descr string
		if len(args) > 5 {
			descr = args[5]
		}
		applyBoth, errBool := strconv.ParseBool(both)
		if errBool != nil {
			log.Printf("FAILURE: parse bool error: %v: apply-both-directions=%v", errBool, both)
			return
		}
		errAdd := a.ContractSubjectAdd(tenant, contract, subject, reverse, applyBoth, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s %s", tenant, contract, subject, reverse, descr)
	case "del":
		if len(args) < 3 {
			log.Fatalf("usage: %s del tenant contract subject", os.Args[0])
		}
		tenant := args[0]
		contract := args[1]
		subject := args[2]
		errDel := a.ContractSubjectDel(tenant, contract, subject)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s %s", tenant, contract, subject)
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
