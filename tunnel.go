//thanks shadowsocks-go
package main

import (
	"log"
	"net"
	"sync"
	"time"
)

func ListenAndServe(addr string, remote func() string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go HandleConnection(c, remote())
	}
}

func SetReadTimeout(c net.Conn) {
	c.SetReadDeadline(time.Now().Add(time.Second * 30))

}

var p = &sync.Pool{New: func() interface{} { return make([]byte, 4096) }}

func PipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	buf := p.Get().([]byte)
	defer p.Put(buf)
	var n int
	var err, err1 error
	for {
		SetReadTimeout(src)
		n, err = src.Read(buf)
		if n > 0 {
			if _, err1 = dst.Write(buf[:n]); err1 != nil {
				break
			}

		}
		if err != nil {
			break
		}
	}
}

func HandleConnection(local net.Conn, remoteAddr string) {
	remote, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer remote.Close()
	go PipeThenClose(local, remote)
	PipeThenClose(remote, local)
}
