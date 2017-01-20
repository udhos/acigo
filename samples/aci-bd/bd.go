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
		log.Fatalf("usage: %s add|del|list|vrf-set|vrf-get|subnet-add|subnet-del|subnet-get|subnet-scope-set|subnet-scope-get|l3out-add|l3out-del args", os.Args[0])
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

	list, errList := a.BridgeDomainList(tenant)
	if errList != nil {
		log.Printf("could not list bridge domains: %v", errList)
		return
	}

	for _, t := range list {
		name := t["name"]
		dn := t["dn"]
		mac := t["mac"]
		descr := t["descr"]
		log.Printf("found bridge domain: name=%s dn=%s mac=%s descr=%s", name, dn, mac, descr)

		bd, isStr := name.(string)
		if !isStr {
			log.Printf("bridge domain name is not string: %s", name)
			continue
		}

		vrf, errVrfGet := a.BridgeDomainVrfGet(tenant, bd)
		if errVrfGet == nil {
			log.Printf("  bridge domain %s vrf=[%s]", bd, vrf)
		}

		subnets, errSubnets := a.BridgeDomainSubnetList(tenant, bd)
		if errSubnets == nil {
			for _, s := range subnets {
				ip := s["ip"]
				sDn := s["dn"]
				scope := s["scope"]
				sDescr := s["descr"]
				log.Printf("  bridge domain %s subnet: ip=%s dn=%s scope=%s descr=%s", bd, ip, sDn, scope, sDescr)
			}
		}

		outs, errOuts := a.BridgeDomainL3ExtOutList(tenant, bd)
		if errOuts == nil {
			for _, o := range outs {
				oName := o["tnL3extOutName"]
				oDn := o["dn"]
				log.Printf("  bridge domain %s L3ExtOut: name=%s dn=%s", bd, oName, oDn)
			}
		}

	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 2 {
			log.Fatalf("usage: %s add tenant bridge-domain [descr]", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		var descr string
		if len(args) > 2 {
			descr = args[2]
		}
		errAdd := a.BridgeDomainAdd(tenant, bd, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s", tenant, bd)
	case "del":
		if len(args) < 2 {
			log.Fatalf("usage: %s del tenant bridge-domain", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		errDel := a.BridgeDomainDel(tenant, bd)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", tenant, bd)
	case "list":
	case "vrf-set":
		if len(args) < 3 {
			log.Fatalf("usage: %s vrf-set tenant bridge-domain vrf", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		vrf := args[2]
		errAdd := a.BridgeDomainVrfSet(tenant, bd, vrf)
		if errAdd != nil {
			log.Printf("FAILURE: vrf-set error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: vrf-set: tenant=%s bd=%s vrf=%s", tenant, bd, vrf)
	case "vrf-get":
		if len(args) < 2 {
			log.Fatalf("usage: %s vrf-get tenant bridge-domain", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		vrf, errGet := a.BridgeDomainVrfGet(tenant, bd)
		if errGet != nil {
			log.Printf("FAILURE: vrf-set error: %v", errGet)
			return
		}
		log.Printf("SUCCESS: vrf-get: tenant=%s bd=%s: => vrf=%s", tenant, bd, vrf)
	case "subnet-add":
		if len(args) < 3 {
			log.Fatalf("usage: %s subnet-add tenant bridge-domain subnet [descr]", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		subnet := args[2]
		var descr string
		if len(args) > 3 {
			descr = args[3]
		}
		errAdd := a.BridgeDomainSubnetAdd(tenant, bd, subnet, descr)
		if errAdd != nil {
			log.Printf("FAILURE: subnet-add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: subnet-add: tenant=%s bd=%s subnet=%s", tenant, bd, subnet)
	case "subnet-del":
		if len(args) < 3 {
			log.Fatalf("usage: %s subnet-del tenant bridge-domain subnet", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		subnet := args[2]
		errDel := a.BridgeDomainSubnetDel(tenant, bd, subnet)
		if errDel != nil {
			log.Printf("FAILURE: subnet-del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: subnet-del: tenant=%s bd=%s subnet=%s", tenant, bd, subnet)
	case "subnet-get":
		if len(args) < 3 {
			log.Fatalf("usage: %s subnet-get tenant bridge-domain subnet", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		subnet := args[2]
		sn, errDel := a.BridgeDomainSubnetGet(tenant, bd, subnet)
		if errDel != nil {
			log.Printf("FAILURE: subnet-get error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: subnet-get: tenant=%s bd=%s subnet=%s: %v", tenant, bd, subnet, sn)
	case "subnet-scope-set":
		if len(args) < 4 {
			log.Fatalf("usage: %s subnet-scope-set tenant bridge-domain subnet scope", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		subnet := args[2]
		scope := args[3]
		errScope := a.BridgeDomainSubnetScopeSet(tenant, bd, subnet, scope)
		if errScope != nil {
			log.Printf("FAILURE: subnet-scope-set error: %v", errScope)
			return
		}
		log.Printf("SUCCESS: subnet-scope-set: tenant=%s bd=%s subnet=%s scope=%s", tenant, bd, subnet, scope)
	case "subnet-scope-get":
		if len(args) < 3 {
			log.Fatalf("usage: %s subnet-scope-get tenant bridge-domain subnet", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		subnet := args[2]
		scope, errScope := a.BridgeDomainSubnetScopeGet(tenant, bd, subnet)
		if errScope != nil {
			log.Printf("FAILURE: subnet-scope-get error: %v", errScope)
			return
		}
		log.Printf("SUCCESS: subnet-scope-get: tenant=%s bd=%s subnet=%s: => scope=%s", tenant, bd, subnet, scope)
	case "l3out-add":
		if len(args) < 3 {
			log.Fatalf("usage: %s l3out-add tenant bridge-domain out", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		out := args[2]
		errAdd := a.BridgeDomainL3ExtOutAdd(tenant, bd, out)
		if errAdd != nil {
			log.Printf("FAILURE: l3out-add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: l3out-add: tenant=%s bd=%s out=%s", tenant, bd, out)
	case "l3out-del":
		if len(args) < 3 {
			log.Fatalf("usage: %s l3out-del tenant bridge-domain out", os.Args[0])
		}
		tenant := args[0]
		bd := args[1]
		out := args[2]
		errDel := a.BridgeDomainL3ExtOutDel(tenant, bd, out)
		if errDel != nil {
			log.Printf("FAILURE: l3out-del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: l3out-del: tenant=%s bd=%s out=%s", tenant, bd, out)
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
