# DNS Library

This project is a personal learning exercise to understand the internals of the DNS protocol and practice writing idiomatic Go. 

To closely familiarise myself with how the protocol works under the hood, I am implementing the binary encoding and decoding completely from scratch using only Go's standard library—without external dependencies or the low-level `golang.org/x/net/dns/dnsmessage` package.

The library aims to support the core wire format specified in RFC 1035, focused on a few fundamental resource record types:
- [ ] `A`
- [ ] `AAAA`
- [ ] `CNAME`
- [ ] `TXT`

