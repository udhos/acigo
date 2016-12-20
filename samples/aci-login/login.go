package main

import (
	"log"

	"github.com/udhos/acigo/aci"
)

func main() {
	a, errNew := aci.New(aci.ClientOptions{Debug: true})
	if errNew != nil {
		log.Printf("login new client error: %v", errNew)
		return
	}

	errLogin := a.Login()
	if errLogin != nil {
		log.Printf("login error: %v", errLogin)
		return
	}

	log.Printf("login ok")
}
