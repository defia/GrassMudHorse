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
	Mailconfig  MailConfig
}
type MailConfig struct {
	User            string
	Password        string
	Host            string
	To              string
	WarningThrottle float64
	RecoverThrottle float64
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
