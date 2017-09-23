package gdns

import (
	"errors"
	"net/http"

	"github.com/miekg/dns"
)

// ResolveRequest parse request data to struct
type ResolveRequest struct {
	Name string // string, required
	// The only required parameter.
	// Its length must be between 1 and 253 (ignoring an optional trailing dot if present).
	// All labels (parts of the name separated by dots) must be 1 to 63 bytes long.
	// The API does not support names with escaped or non-ASCII characters, but they are not explicitly rejected.
	// Internationalized domain names must use punycode format (e.g. "xn--qxam" rather than "ελ").
	Type uint16 // string, default: 1
	// RR type can be represented as a number in [1, 65535] or a canonical string (case-insensitive, such as A or aaaa).
	// You can use 255 for 'ANY' queries but be aware that this is not a replacement for sending queries for both
	// A and AAAA or MX records. Authoritative name servers need not return all records for such queries;
	// some do not respond, and others (such as cloudflare.com) return only HINFO.
	CheckingDisabled bool // boolean, default: false
	// The CD (checking disabled) bit. Use cd, cd=1, or cd=true to disable DNSSEC validation;
	// use cd=0, cd=false, or no cd parameter to enable DNSSEC validation.
	EDNSClientSubnet string // string, default: empty
	// The edns0-client-subnet option. Format is an IP address with a subnet mask. Examples: 1.2.3.4/24, 2001:700:300::/48.
	// If you are using DNS-over-HTTPS because of privacy concerns,
	// and do not want any part of your IP address to be sent to authoritative name servers for geographic location accuracy,
	// use edns_client_subnet=0.0.0.0/0. Google Public DNS normally sends approximate network information
	// (usually zeroing out the last part of your IPv4 address).
	RandomPadding string // string, ignored
	// The value of this parameter is ignored. Example: XmkMw~o_mgP2pf.gpw-Oi5dK.
	// API clients concerned about possible side-channel privacy attacks
	// using the packet sizes of HTTPS GET requests can use this to make all requests
	// exactly the same size by padding requests with random data.
	// To prevent misinterpretation of the URL, restrict the padding characters to
	// the unreserved URL characters:
	// upper- and lower-case letters, digits, hyphen, period, underscore and tilde.
}

// Request Google DNS API
func (r *ResolveRequest) Request() (response []byte, statusCode int, err error) {
	if r.Name == "" {
		err = errors.New("invalid resolve name")
		return
	}

	if r.Type == dns.TypeNone {
		r.Type = dns.TypeANY
	}

	// make request URL params
	params := map[string]interface{}{
		"name": r.Name,
		"type": r.Type,
		"cd":   r.CheckingDisabled,
	}

	if r.EDNSClientSubnet != "" {
		params["edns_client_subnet"] = r.EDNSClientSubnet
	}

	response, statusCode, err = httpGet(http.MethodGet, GoogleDNSAPI+"resolve", params)
	return
}
