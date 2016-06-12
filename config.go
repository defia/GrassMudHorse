package main

import (
	"encoding/json"
	"os"
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
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = json.NewDecoder(f).Decode(config)
	return config, err
}
