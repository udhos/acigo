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

	// add/del nodes

	execute(a, os.Args[1], os.Args[2:])

	// display existing nodes
	nodes, errList := a.NodeList()
	if errList != nil {
		log.Printf("could not list nodes: %v", errList)
		return
	}

	for _, n := range nodes {
		name := n["name"]
		dn := n["dn"]
		role := n["role"]
		fmt.Printf("FOUND node: name=%s role=%s dn=%s\n", name, role, dn)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 3 {
			log.Fatalf("usage: %s add name ID serial", os.Args[0])
		}
		name := args[0]
		ID := args[1]
		serial := args[2]
		errAdd := a.NodeAdd(name, ID, serial)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add node: %s", name)
	case "del":
		if len(args) < 1 {
			log.Fatalf("usage: %s del serial", os.Args[0])
		}
		serial := args[0]
		errDel := a.NodeDel(serial)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del node: %s", serial)
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

	// Since credentials have not been specified explicitly under ClientOptions,
	// Login() will use env vars: APIC_HOSTS=host, APIC_USER=username, APIC_PASS=pwd
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

	log.Printf("logout: done")
}
