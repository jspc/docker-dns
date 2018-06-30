package main

import (
	"flag"
	"log"

	"github.com/miekg/dns"
)

var (
	domain   = flag.String("domain", "internal.", "Domain name with which to serve records on")
	target   = flag.String("target", "localhost", "Address (host, ipv4, ipv6 etc.) to respond to requests with")
	resolver = flag.String("resolver", "1.1.1.1:53", "DNS Server to pass reqeusts to which are not on our domain")
	listen   = flag.String("listen", ":53", "Port on which to serve DNS requests")
)

func main() {
	flag.Parse()

	log.Print("Starting")

	docker, err := NewDocker()
	if err != nil {
		log.Panic(err)
	}

	d, err := NewDNS(*domain, *target, *resolver, docker)
	if err != nil {
		log.Panic(err)
	}

	panic(dns.ListenAndServe(*listen, "udp", d))
}
