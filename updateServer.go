// Copyright 2014 Markus Piéton (marpie). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	RET_SUCCESSFUL   = "good"
	RET_NO_CHANGE    = "nochg"
	RET_INVALID_AUTH = "badauth"
	RET_INVALID_HOST = "nohost"
	RET_NO_HOST_NAME = "notfqdn"
)

func HttpUpdateHandler(w http.ResponseWriter, req *http.Request) {
	// Get client IP address
	client_ip := strings.Split(req.RemoteAddr, ":")[0]
	if req.Header.Get("Cf-Connecting-Ip") != "" {
		client_ip = req.Header.Get("Cf-Connecting-Ip")
	} else if req.Header.Get("X-Forwarded-For") != "" {
		client_ip = req.Header.Get("X-Forwarded-For")
	}

	if req.ParseForm() != nil {
		return
	}

	hostname := GetIgnoreCase(req.Form, "hostname")
	if hostname == "" {
		fmt.Fprint(w, RET_NO_HOST_NAME)
		return
	}

	// Retrieve login details
	username, password := DecodeBasicAuth(req.Header.Get("Authorization"))
	// Check credentials
	if (username == "") && (password == "") {
		// Output help
		fmt.Fprintln(w, "/nic/update?hostname=<domain>&myip=<ipaddr>")
		return
	}

	user, err := config.GetUser(username, password)
	if err != nil {
		// Invalid Login
		fmt.Fprint(w, RET_INVALID_AUTH)
		return
	}

	if !user.HasDomain(hostname) {
		fmt.Fprint(w, RET_INVALID_HOST)
		return
	}

	myip := GetIgnoreCase(req.Form, "myip")
	if myip != "" {
		client_ip = myip
	}

	// Sanitize IP Address
	if ip := net.ParseIP(client_ip); ip.IsUnspecified() {
		fmt.Fprint(w, RET_NO_CHANGE)
		return
	} else {
		client_ip = ip.String()
	}

	updHandler.Add(hostname, client_ip)

	fmt.Fprint(w, RET_SUCCESSFUL)
}
