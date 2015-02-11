package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Probe struct {
	Servers       []*Server
	Interval      time.Duration
	fastestServer *Server
}

type Server struct {
	Address string
	P       *pinger
	Stat    *Stat
}

func (s *Server) String() string {
	return fmt.Sprintf("server:%s stats in recent %d pings:average latency:%s package drop rate:%f%%", s.P.Address, s.Stat.Count, s.Stat.AverageLatency().String(), s.Stat.DropRate()*100)
}

func (p *Probe) GetFastestServer() string {
	return p.fastestServer.Address
}
func (p *Probe) getFastestServer() *Server {
	max := p.Servers[0]
	maxscore := p.Servers[0].Stat.Score()
	for _, v := range p.Servers {
		debugOut(v)
		if score := v.Stat.Score(); score > maxscore {
			maxscore = score
			max = v
		}
	}
	debugOut("max:", max.Address)
	return max
}

func NewProbe(config *Config) *Probe {
	p := new(Probe)
	length := len(config.Servers)
	if length <= 0 {
		log.Fatal("server list empty")
	}

	p.Servers = make([]*Server, length)
	var version ICMPVersion
	if config.IPv6 {
		version = ICMPv6
	} else {
		version = ICMPv4
	}
	p.Interval = time.Duration(config.Interval) * time.Millisecond
	for i, v := range config.Servers {
		addr := CutPort(v)

		p.Servers[i] = &Server{
			Address: v,
			P:       NewPinger(addr, config.Timeout, config.PayloadSize, version),
			Stat:    NewStat(config.HistorySize),
		}
	}
	p.fastestServer = p.Servers[0]
	return p
}

func CutPort(addr string) string {
	return strings.Split(addr, ":")[0]
}

func (p *Probe) Start() {
	go func() {
		for {
			p.fastestServer = p.getFastestServer()
			time.Sleep(time.Second)
		}
	}()
	for _, v := range p.Servers {
		go func(v *Server) {
			ticker := time.Tick(p.Interval)
			for _ = range ticker {
				go func() {
					latency, err := v.P.Ping()
					v.Stat.AddResult(latency, err != nil)
				}()
			}
		}(v)
	}
}
