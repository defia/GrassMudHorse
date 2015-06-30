package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Listen      string
	Interval    int
	Timeout     int
	IPv6        bool
	PayloadSize int
	HistorySize int
	Servers     []string
	Lua         string
}

func ReadConfig(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = json.Unmarshal(b, config)
	return config, err
}
