package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Probe struct {
	Servers            []*Server
	Interval           time.Duration
	fastestServer      *Server
	config             *Config
	recentsendmailtime time.Time
	badflag            bool
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
func (p *Probe) getFastestServer() (*Server, float64) {
	max := p.Servers[0]
	maxscore := p.Servers[0].Stat.Score()
	mindroprate := p.Servers[0].Stat.DropRate()
	for _, v := range p.Servers {
		debugOut(v)
		if score := v.Stat.Score(); score > maxscore {
			maxscore = score
			max = v
		}
		if droprate := v.Stat.DropRate(); droprate < mindroprate {
			mindroprate = droprate
		}
	}
	debugOut("max:", max.Address)
	return max, mindroprate
}

func NewProbe(config *Config) *Probe {

	p := new(Probe)
	p.badflag = false
	p.config = config
	p.recentsendmailtime = time.Unix(0, 0)
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

		mailconfig := p.config.Mailconfig
		var rate float64
		for {
			p.fastestServer, rate = p.getFastestServer()
			if !p.badflag && rate > mailconfig.WarningThrottle && time.Now().Sub(p.recentsendmailtime) > time.Minute*5 {
				log.Println("bad")
				err := SendMail(mailconfig.User, mailconfig.Password, mailconfig.Host, mailconfig.To, true)
				if err == nil {

					p.badflag = true
					p.recentsendmailtime = time.Now()

				} else {
					log.Println(err)
				}
				continue
			}
			if p.badflag && mailconfig.RecoverThrottle > rate {
				log.Println("recover")
				err := SendMail(mailconfig.User, mailconfig.Password, mailconfig.Host, mailconfig.To, false)
				if err == nil {
					p.badflag = false
					p.recentsendmailtime = time.Now()

				} else {
					log.Println(err)
				}

			}

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
