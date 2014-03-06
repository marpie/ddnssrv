// Copyright 2014 Markus Piéton (marpie). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// command ddnssrv is a simple dynamic DNS listener that accepts
// DNS updates in the same format as the DynDNS.org servers.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	config     *Config
	updHandler *UpdateHandler
)

func AliveHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Alive!")
}

// DefaultHandler is called if no other handler matches.
func DefaultHandler(w http.ResponseWriter, req *http.Request) {
	// Do nothing and close connection.
	return
}

func main() {
	var err error
	var ip = flag.String("ip", "", "IP address to listen on.")
	var port = flag.Uint("port", 8080, "Port to listen on.")
	var mgmt_port = flag.Uint("mgmt-port", 8090, "Port to listen on.")
	var configFile = flag.String("cfg", "config.json", "Configuration file to use.")
	var createConfig = flag.Bool("cfg-create", false, "If set a default config is written.")
	flag.Parse()

	if *createConfig {
		// Create default config.
		fmt.Printf("Writing config file: %v\n", *configFile)
		var cfg Config

		usr, err := cfg.AddUser("Default-User", "MyPassword")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
			os.Exit(1)
		}
		usr.AddDomain("www.example.com")

		if err := cfg.Save(*configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Writing config: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	config, err = LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// Create Update Handler to syncronize MgmtServer and FrontendServer
	updHandler = NewUpdateHandler()

	// ----------------------------------------------------------------------
	// Initialize Management Server
	StartMgmtSrv(*mgmt_port)

	// ----------------------------------------------------------------------
	// Initialize Web Server
	http.HandleFunc("/status", AliveHandler)
	http.HandleFunc("/nic/update", HttpUpdateHandler)
	http.HandleFunc("/", DefaultHandler)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *ip, *port), nil); err != nil {
		fmt.Fprintf(os.Stderr, "[DnsUpdateSrv] ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
