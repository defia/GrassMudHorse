package main

import (
	"net"
	"os"
	"sync"
	"time"
)

var (
	globalseq int32 = 0
	lock            = &sync.Mutex{}
	xid             = os.Getpid() & 0xffff
)

type pinger struct {
	Address      string
	Timeout      time.Duration
	Data         []byte
	Type         ICMPVersion
	responceType int
	requestType  int
	dialNetwork  string
	conn         net.Conn
}

type ICMPVersion int

const (
	ICMPv6 ICMPVersion = 6
	ICMPv4 ICMPVersion = 4
)

func NewPinger(address string, timeoutMillisecond int, size int, typ ICMPVersion) *pinger {

	p := pinger{
		Address: address,
		Timeout: time.Millisecond * time.Duration(timeoutMillisecond),
		Data:    make([]byte, size),
		Type:    typ,
	}

	if typ == ICMPv6 {
		p.requestType = icmpv6EchoRequest
		p.responceType = icmpv6EchoReply
		p.dialNetwork = "ip6:ipv6-icmp"
	} else {

		p.requestType = icmpv4EchoRequest
		p.responceType = icmpv4EchoReply
		p.dialNetwork = "ip4:icmp"

	}

	return &p
}

func (p *pinger) Ping() (latency time.Duration, err error) {
	c, err := net.Dial(p.dialNetwork, p.Address)
	if err != nil {
		return
	}
	defer c.Close()

	lock.Lock()
	seq := int(globalseq)
	globalseq++
	if globalseq > 65535 {
		globalseq = 0
	}
	lock.Unlock()

	wb, err := (&icmpMessage{
		Type: p.requestType, Code: 0,
		Body: &icmpEcho{
			ID: xid, Seq: seq,
			Data: p.Data,
		},
	}).Marshal()
	if err != nil {
		return
	}
	t := time.Now()
	c.SetDeadline(t.Add(p.Timeout))
	if _, err = c.Write(wb); err != nil {
		return
	}
	var m *icmpMessage
	rb := make([]byte, 20+len(wb))
	for {
		if _, err = c.Read(rb); err != nil {
			return
		}
		if p.Type == ICMPv4 {
			rb = ipv4Payload(rb)
		}
		if m, err = parseICMPMessage(rb); err != nil {
			return
		}
		switch m.Type {
		case icmpv4EchoReply, icmpv6EchoReply:
			var b []byte
			if b, err = m.Body.Marshal(); err != nil {
				return
			} else {
				var echo *icmpEcho
				//TODO add icmp checksum verify, see http://tools.ietf.org/html/rfc1071
				if echo, err = parseICMPEcho(b); err != nil {
					return
				} else {
					//check if type mismatch or packet sequence/id mismatch
					if echo.Seq != seq || echo.ID != xid || m.Type != p.responceType {
						continue
					}
					return time.Now().Sub(t), nil
				}

			}
			break
		}
	}
	return

}
