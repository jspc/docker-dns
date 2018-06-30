package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type DNSClient interface {
	Dial(string) (*dns.Conn, error)
	Exchange(*dns.Msg, string) (*dns.Msg, time.Duration, error)
}

type DNS struct {
	Client   DNSClient
	Docker   Docker
	Domain   string
	NextHost string

	Target        *net.IPAddr
	TargetAddress string
	TargetType    string
}

func NewDNS(domain, targetAddress, serverAddress string, docker Docker) (d DNS, err error) {
	d.Client = new(dns.Client)
	d.Docker = docker
	d.Domain = domain
	d.NextHost = serverAddress
	d.TargetAddress = targetAddress

	d.Target, err = net.ResolveIPAddr("ip", targetAddress)
	if err != nil {
		return
	}

	switch {
	case d.Target.String() != targetAddress:
		if !strings.HasSuffix(d.TargetAddress, ".") {
			d.TargetAddress = fmt.Sprintf("%s.", d.TargetAddress)
		}

		d.TargetType = "CNAME"
	case d.Target.IP.To4() != nil:
		d.TargetType = "A"
	default:
		d.TargetType = "AAAA"
	}

	return
}

func (d DNS) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	question := r.Question[0].Name
	components := strings.SplitN(question, ".", 2)

	if len(components) != 2 || components[1] != d.Domain {
		m, _, err := d.Client.Exchange(r, d.NextHost)
		if err != nil {
			log.Print(err)
		}

		w.WriteMsg(m)

		return
	}

	containers, err := d.Docker.Containers()
	if err != nil {
		log.Print(err)

		w.WriteMsg(m)

		return
	}

	host := components[0]
	_, ok := containers[host]
	if ok {
		m.Answer = []dns.RR{d.CreateRecord(question)}
	}

	w.WriteMsg(m)
}

func (d DNS) CreateRecord(h string) dns.RR {
	switch d.TargetType {
	case "CNAME":
		r := dns.CNAME{}

		r.Target = d.TargetAddress
		r.Hdr = dns.RR_Header{Name: h, Class: 1, Rrtype: dns.TypeCNAME, Ttl: 0}

		return &r

	case "A":
		r := dns.A{}

		r.A = d.Target.IP
		r.Hdr = dns.RR_Header{Name: h, Class: 1, Rrtype: dns.TypeA, Ttl: 0}

		return &r
	}

	r := dns.AAAA{}

	r.AAAA = d.Target.IP
	r.Hdr = dns.RR_Header{Name: h, Class: 1, Rrtype: dns.TypeAAAA, Ttl: 0}

	return &r
}
