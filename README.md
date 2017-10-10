# ΛPΞX [WIP]
[![Build Status](https://travis-ci.org/oif/apex.svg?branch=master)](https://travis-ci.org/oif/apex)
[![Coverage Status](https://coveralls.io/repos/github/oif/apex/badge.svg?branch=master)](https://coveralls.io/github/oif/apex?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/oif/apex)](https://goreportcard.com/report/github.com/oif/apex)

ΛPΞX is a DNS server written in Go, help you connect to the real Internet.

# Feature
* edns-client-subnet
* Multiple DNS upstream
    * DNS Over TLS
    * UDP/TCP DNS
* Hosts
* DNSSEC
* Cache

# Package Management
Apex uses the Go community [dep](https://github.com/golang/dep) project for package management, but it's young so maybe will cause some unexcepted issue occurred during building. For more information about dep project status, check [dep - Current status](https://github.com/golang/dep#current-status)

# License
MIT
