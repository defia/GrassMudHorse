package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
)

const FileName = "config.json"

var debug *bool

func main() {
	debug = flag.Bool("debug", false, "output debug msg")
	web := flag.String("web", ":1234", "web log listening ip:port")
	config, err := ReadConfig(FileName)
	flag.Parse()
	if err != nil {
		log.Fatal(err)
	}
	p := NewProbe(config)
	go p.Start()
	go ListenAndServe(config.Listen, p.GetFastestServer)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body>")
		for _, v := range p.Servers {

			fmt.Fprintf(w, "%s<br>", html.EscapeString(v.String()))
		}
		fmt.Fprintf(w, "now use:%s", p.GetFastestServer())
	})
	log.Fatal(http.ListenAndServe(*web, nil))
}
func debugOut(d ...interface{}) {
	if *debug {
		log.Println(d)
	}
}
