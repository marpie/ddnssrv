// Copyright 2014 Markus Piéton (marpie). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
)

type Update struct {
	Hostname    string
	Destination string
}

type UpdateHandler struct {
	updates    map[string]string
	quitChan   chan bool
	newUpdate  chan *Update
	getUpdates chan chan string
}

func NewUpdateHandler() *UpdateHandler {
	handler := UpdateHandler{}
	handler.updates = make(map[string]string)
	handler.quitChan = make(chan bool, 1)
	handler.newUpdate = make(chan *Update, 1)
	handler.getUpdates = make(chan chan string, 1)

	go handler.handle()

	return &handler
}

func (handler *UpdateHandler) Add(hostname, destination string) {
	handler.newUpdate <- &Update{hostname, destination}
}

func (handler *UpdateHandler) GetUpdatesAsString() string {
	upds := make(chan string, 1)
	handler.getUpdates <- upds
	return <-upds
}

func (handler *UpdateHandler) Stop() {
	handler.quitChan <- true
}

func (handler *UpdateHandler) handle() {
	for {
		select {
		case getUpdates := <-handler.getUpdates:
			res := make([]string, len(handler.updates))
			i := 0
			for hostname, destination := range handler.updates {
				res[i] = hostname + "=" + destination
			}
			getUpdates <- strings.Join(res, ";")
			// Remove all updates
			handler.updates = make(map[string]string)
		case update := <-handler.newUpdate:
			handler.updates[update.Hostname] = update.Destination
		case <-handler.quitChan:
			return
		}
	}
}
