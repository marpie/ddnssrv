// Copyright 2014 Markus Piéton (marpie). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

var (
	ErrInvalidUser = errors.New("The username or password is invalid.")
)

type User struct {
	Username string   `json:"Username,omitempty"`
	Password string   `json:"Password,omitempty"`
	Domains  []string `json:"Domains"`
}

func (usr *User) AddDomain(domain string) {
	if usr.HasDomain(domain) {
		return
	}
	usr.Domains = append(usr.Domains, domain)
}

func (usr *User) HasDomain(domain string) bool {
	domain = strings.ToLower(domain)
	for _, entry := range usr.Domains {
		if strings.ToLower(entry) == domain {
			return true
		}
	}
	return false
}

type Config struct {
	Users []User `json:"Users,omitempty"`
}

func NewConfig() (config *Config) {
	config = &Config{}
	config.Users = make([]User, 0)

	return config
}

func LoadConfig(filename string) (config *Config, err error) {
	// Read config
	config_json, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	// Parse config
	config = &Config{}
	err = json.Unmarshal(config_json, config)
	if err != nil {
		return
	}

	return
}

func (cfg *Config) AddUser(username, password string) (*User, error) {
	// Check if user already exists
	usr, err := cfg.GetUser(username, password)
	if err == nil {
		return usr, nil
	}

	user := User{}
	user.Username = username
	user.Password = password
	user.Domains = make([]string, 0)
	cfg.Users = append(cfg.Users, user)

	return cfg.GetUser(username, password)
}

func (cfg *Config) GetUser(username, password string) (*User, error) {
	for idx, user := range cfg.Users {
		if (user.Username == username) && (user.Password == password) {
			return &cfg.Users[idx], nil
		}
	}
	return nil, ErrInvalidUser
}

func (cfg *Config) Save(filename string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}
