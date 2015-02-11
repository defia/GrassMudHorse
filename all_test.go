package main

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_Pinger(t *testing.T) {
	p := NewPinger("baidu.com", 500, 32, ICMPv4)

	if latency, err := p.Ping(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(latency)
	}

	p = NewPinger("::1", 500, 32, ICMPv6)

	if latency, err := p.Ping(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(latency)
	}
}

var sampleconfig = &Config{
	Listen:      ":1080",
	Interval:    1000,
	Timeout:     500,
	IPv6:        false,
	PayloadSize: 1024,
	HistorySize: 100,
	Servers:     []string{"8.8.8.8:5678", "192.168.1.1:1234"},
}

func Test_ReadConfig(t *testing.T) {

	c, err := ReadConfig("sample.json")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, sampleconfig) {
		t.Fatal(c, "shoule be ", sampleconfig)
	}
}

func Test_Stat(t *testing.T) {
	latency := []time.Duration{0, 200, 150, 0}
	drop := []bool{true, false, false, true}
	averagelatency := []time.Duration{time.Second * 10, 200, 175, 150}
	droprate := []float64{1, 0.5, 0, 0.5}
	stat := NewStat(2)
	for i := 0; i < 4; i++ {
		stat.AddResult(latency[i], drop[i])
		if stat.AverageLatency() != averagelatency[i] || stat.DropRate() != droprate[i] {
			log.Fatal(stat.AverageLatency(), averagelatency[i], stat.DropRate(), droprate[i])
		}
	}

}

func Test_GetFastestServer(t *testing.T) {

}
