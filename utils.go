// Copyright 2014 Markus Piéton (marpie). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func DebugRequest(req *http.Request) {
	fmt.Printf("\nNew Request: %s\n", req.RemoteAddr)
	fmt.Printf("%v %v %v\n", req.Method, req.RequestURI, req.Proto)
	for key, values := range req.Header {
		fmt.Printf("%s: %s\n", key, strings.Join(values, ";"))
	}
}

func DecodeBasicAuth(basicAuth string) (username, password string) {
	parts := strings.Split(basicAuth, " ")
	if (len(parts) < 2) || (parts[0] != "Basic") {
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return
	}
	parts = strings.Split(string(decoded), ":")
	if len(parts) < 2 {
		return
	}
	return parts[0], parts[1]
}

func GetIgnoreCase(values url.Values, key string) string {
	key = strings.ToLower(key)
	for entry_key, values := range values {
		if strings.ToLower(entry_key) == key {
			if len(values) < 1 {
				return ""
			}
			return values[0]
		}
	}
	return ""
}
