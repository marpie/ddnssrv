// Copyright 2014 Markus Piéton (marpie). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"os"
)

func GetUpdatesHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, updHandler.GetUpdatesAsString())
}

func StartMgmtSrv(port uint) {
	mgmt_srv := http.NewServeMux()
	mgmt_srv.HandleFunc("/status", AliveHandler)
	mgmt_srv.HandleFunc("/get/updates", GetUpdatesHandler)
	mgmt_srv.HandleFunc("/", DefaultHandler)

	go func(mgmt_srv *http.ServeMux, port uint) {

		srv := &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", port),
			Handler: mgmt_srv,
		}

		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "[MgmtInterface] ListenAndServe: %v\n", err)
			os.Exit(1)
		}
	}(mgmt_srv, port)
}
