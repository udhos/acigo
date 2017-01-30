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

	list, errList := a.RemoteLocationList()
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}

	for _, t := range list {
		name := t["name"]
		host := t["host"]
		protocol := t["protocol"]
		remotePort := t["remotePort"]
		remotePath := t["remotePath"]
		username := t["userName"]
		descr := t["descr"]

		log.Printf("FOUND remote location: name=%s host=%s proto=%s remPort=%s remPath=%s user=%s descr=%s", name, host, protocol, remotePort, remotePath, username, descr)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 7 {
			log.Fatalf("usage: %s add location host protocol remotePort remotePath username password [descr]", os.Args[0])
		}
		location := args[0]
		host := args[1]
		protocol := args[2]
		remotePort := args[3]
		remotePath := args[4]
		username := args[5]
		password := args[6]
		var descr string
		if len(args) > 7 {
			descr = args[7]
		}
		errAdd := a.RemoteLocationAdd(location, host, protocol, remotePort, remotePath, username, password, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s %s %s %s %s", location, host, protocol, remotePort, remotePath, username, password, descr)
	case "del":
		if len(args) < 1 {
			log.Fatalf("usage: %s del location", os.Args[0])
		}
		location := args[0]
		errDel := a.RemoteLocationDel(location)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s", location)
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
