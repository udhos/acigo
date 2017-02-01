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
		log.Fatalf("usage: %s add|del|list|run args", os.Args[0])
	}

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	list, errList := a.ExportConfigurationList()
	if errList != nil {
		log.Printf("could not list: %v", errList)
		return
	}

	for _, t := range list {
		config := t["name"]
		dn := t["dn"]
		adminSt := t["adminSt"]
		format := t["format"]
		descr := t["descr"]

		log.Printf("FOUND export config: config=%s dn=%s adminSt=%s format=%s descr=%s", config, dn, adminSt, format, descr)

		conf, isStr := config.(string)
		if !isStr {
			log.Printf("  config=%s not a string", config)
			continue
		}

		loc, errLoc := a.ExportConfigurationRemoteLocationGet(conf)
		if errLoc == nil {
			name := loc["tnFileRemotePathName"]
			log.Printf("  config=%s remote location: name=[%s]", conf, name)
		}

		sched, errSched := a.ExportConfigurationSchedulerGet(conf)
		if errSched == nil {
			name := sched["tnTrigSchedPName"]
			log.Printf("  config=%s scheduler: name=[%s]", conf, name)
		}
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 3 {
			log.Fatalf("usage: %s add config scheduler remote-location [descr]", os.Args[0])
		}
		config := args[0]
		scheduler := args[1]
		remoteLocation := args[2]
		var descr string
		if len(args) > 3 {
			descr = args[3]
		}
		errAdd := a.ExportConfigurationAdd(config, scheduler, remoteLocation, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s %s %s", config, scheduler, remoteLocation, descr)
	case "del":
		if len(args) < 1 {
			log.Fatalf("usage: %s del config", os.Args[0])
		}
		config := args[0]
		errDel := a.ExportConfigurationDel(config)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s", config)
	case "run":
		if len(args) < 1 {
			log.Fatalf("usage: %s run config", os.Args[0])
		}
		config := args[0]
		errRun := a.ExportConfigurationRun(config)
		if errRun != nil {
			log.Printf("FAILURE: run error: %v", errRun)
			return
		}
		log.Printf("SUCCESS: rn: %s", config)
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
