package main

import (
	"github.com/henyxia/screenator/cmd/client"
	"github.com/henyxia/screenator/cmd/server"
	"log"
	"os"
)

func usage() {
	log.Fatalln("usage ./screenator client|server CONFIG_FILE")
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	cmd := os.Args[1]
	configFile := os.Args[2]
	if cmd == "client" {
		client.Client(configFile)
	} else if cmd == "server" {
		server.Server(configFile)
	} else {
		log.Fatalln("retard alert! must choose between clietn and server!")
	}
}
