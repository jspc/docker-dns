Docker DNS
==

Given access to a Docker API and a host on which the containers run on, provide DNS entries for those containers.

Problem Statement
--

1. As a developer, you don't want to have to keep DNS records for containers up to date.
1. As a developer, you have http services which listen on a hostname
1. As a developer, you want to run a proxy service which routes based on requests, but all of your requests go to localhost when running locally

In any of these apply then this project might be useful.

Usage
--

```
 ./docker-dns -h
Usage of ./docker-dns:
  -domain string
        Domain name with which to serve records on (default "internal.")
  -resolver string
        DNS Server to pass reqeusts to which are not on our domain (default "1.1.1.1:53")
  -target string
        Address (host, ipv4, ipv6 etc.) to respond to requests with (default "localhost")
```

The default options will:

 * Receive a DNS query
 * If this query is _not_ on the domain `.internal` then ask `1.1.1.1` to resolve the resquest
   * Otherwise determine whether there's a container which matches the hostname in the request
   * If there is, respond with an answer signifying `target` is the host
   * Otherwise return an empty response

This, on the happy path, looks like:

```bash
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
14adf764ef9d        nginx               "nginx -g 'daemon of…"   3 seconds ago       Up 2 seconds        80/tcp                   my-awesome-site
$ dig @localhost -p 3553 my-awesome-site.internal

; <<>> DiG 9.10.6 <<>> @localhost -p 3553 my-awesome-site.internal
; (2 servers found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 50516
;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;my-awesome-site.internal.      IN      A

;; ANSWER SECTION:
my-awesome-site.internal. 0     IN      CNAME   localhost.

;; Query time: 5 msec
;; SERVER: 127.0.0.1#3553(127.0.0.1)
;; WHEN: Sat Jun 30 16:26:40 BST 2018
;; MSG SIZE  rcvd: 89
```
