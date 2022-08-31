package main

import (
	"flag"
	"sneakyicmp/sneakyicmp"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var body = &icmp.Echo{
	ID:   123123,
	Seq:  100,
	Data: []byte("123213123"),
}

var xd = icmp.Message{
	Type: ipv4.ICMPTypeEcho,
	Code: 0,
	Body: body,
}

func main() {
	var msg string
	var mode bool
	flag.StringVar(&msg, "m", "sneakyleaky", "smt easy")
	flag.BoolVar(&mode, "u", false, "smt easy")

	flag.Parse()
	if !mode {
	} else {
		sneakyicmp.RecvSneakyMessage("172.31.252.17")
	}
}
