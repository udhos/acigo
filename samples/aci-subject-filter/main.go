package main

import (
	"fmt"
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 5 {
		log.Fatalf("usage: %s both-add|both-del|input-add|input-del|output-add|output-del|list args", os.Args[0])
	}

	tenant := os.Args[2]
	contract := os.Args[3]
	subject := os.Args[4]

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	{
		list, errList := a.SubjectFilterBothList(tenant, contract, subject)
		if errList != nil {
			log.Printf("could not list: %v", errList)
			return
		}
		for _, t := range list {
			dn := t["dn"]
			log.Printf("FOUND subject filter both: dn=%s", dn)
		}
	}

	{
		list, errList := a.SubjectFilterInputList(tenant, contract, subject)
		if errList != nil {
			log.Printf("could not list: %v", errList)
			return
		}
		for _, t := range list {
			dn := t["dn"]
			log.Printf("FOUND subject filter input: dn=%s", dn)
		}
	}

	{
		list, errList := a.SubjectFilterOutputList(tenant, contract, subject)
		if errList != nil {
			log.Printf("could not list: %v", errList)
			return
		}
		for _, t := range list {
			dn := t["dn"]
			log.Printf("FOUND subject filter output: dn=%s", dn)
		}
	}

}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "both-add":
		if len(args) < 4 {
			log.Fatalf("usage: %s both-add tenant contract subject filter", os.Args[0])
		}
		tenant := args[0]
		contract := args[1]
		subject := args[2]
		filter := args[3]
		errAdd := a.SubjectFilterBothAdd(tenant, contract, subject, filter)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s", tenant, contract, subject, filter)
	case "both-del":
		if len(args) < 4 {
			log.Fatalf("usage: %s both-del tenant contract subject filter", os.Args[0])
		}
		tenant := args[0]
		contract := args[1]
		subject := args[2]
		filter := args[3]
		errDel := a.SubjectFilterBothDel(tenant, contract, subject, filter)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s %s %s", tenant, contract, subject, filter)
	case "input-add":
		if len(args) < 4 {
			log.Fatalf("usage: %s input-add tenant contract subject filter", os.Args[0])
		}
		tenant := args[0]
		contract := args[1]
		subject := args[2]
		filter := args[3]
		errAdd := a.SubjectFilterInputAdd(tenant, contract, subject, filter)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s", tenant, contract, subject, filter)
	case "input-del":
		if len(args) < 4 {
			log.Fatalf("usage: %s input-del tenant contract subject filter", os.Args[0])
		}
		tenant := args[0]
		contract := args[1]
		subject := args[2]
		filter := args[3]
		errDel := a.SubjectFilterInputDel(tenant, contract, subject, filter)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s %s %s", tenant, contract, subject, filter)
	case "output-add":
		if len(args) < 4 {
			log.Fatalf("usage: %s output-add tenant contract subject filter", os.Args[0])
		}
		tenant := args[0]
		contract := args[1]
		subject := args[2]
		filter := args[3]
		errAdd := a.SubjectFilterOutputAdd(tenant, contract, subject, filter)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s", tenant, contract, subject, filter)
	case "output-del":
		if len(args) < 4 {
			log.Fatalf("usage: %s output-del tenant contract subject filter", os.Args[0])
		}
		tenant := args[0]
		contract := args[1]
		subject := args[2]
		filter := args[3]
		errDel := a.SubjectFilterOutputDel(tenant, contract, subject, filter)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s %s %s", tenant, contract, subject, filter)
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
