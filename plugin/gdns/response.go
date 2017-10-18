package gdns

import (
	"encoding/json"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

// ResolveResponse from Google DNS API
type ResolveResponse struct {
	Status             int           `json:"Status"`             // 0 success, 2 error. Standard DNS response code (32 bit integer)
	Truncated          bool          `json:"TC"`                 // Whether the response is truncated
	RecursionDesired   bool          `json:"RD"`                 // RD, Always true for Google Public DNS
	RecursionAvailable bool          `json:"RA"`                 // RA, Always true for Google Public DNS
	AuthenticatedData  bool          `json:"AD"`                 // AD, Whether all response data was validated with DNSSEC
	CheckingDisabled   bool          `json:"CD"`                 // CD, Whether the client asked to disable DNSSEC
	Question           []Question    `json:"Question"`           // Question
	Answer             []Answer      `json:"Answer"`             // Answer
	Authority          []Answer      `json:"Authority"`          // Authority
	Additional         []interface{} `json:"Additional"`         // Additional response
	EDNSClientSubnet   string        `json:"edns_client_subnet"` // IP address / scope prefix-length, ref. https://tools.ietf.org/html/draft-ietf-dnsop-edns-client-subnet-08#section-6
	Comment            string        `json:"Comment"`            // Comment
}

// Question part of response
type Question struct {
	Name string `json:"name"` // FQDN with trailing dot
	Type uint16 `json:"type"` // Standard DNS RR type
}

// Answer part of response
type Answer struct {
	Name string `json:"name"` // Always matches name in the Question section
	Type uint16 `json:"type"` // Standard DNS RR type
	TTL  uint32 `json:"TTL"`  // Record's time-to-live in seconds
	Data string `json:"data"` // IP address as text
}

// BytesToResolveResponse convert resolve request response(bytes) to ResolveResponse
func BytesToResolveResponse(bytes []byte) (rr *ResolveResponse, err error) {
	rr = new(ResolveResponse)
	err = json.Unmarshal(bytes, rr)
	return
}

// Success if response status is 0 means success, otherwise will have a comment for failure detail
func (r *ResolveResponse) Success() (bool, string) {
	return r.Status != 1 && r.Status != 2 && r.Status != 5, r.Comment
}

// ToRR convert a google dns anwser to dns.RR
// @TODO refactor useing reflect feature
func (a Answer) ToRR() (rr dns.RR) {
	// Currently support part of rr type only
	// reuse code from github.com/oif/proton/gdns/response.go
	switch a.Type {
	case dns.TypeA:
		rr = &dns.A{
			Hdr: a.GetRRHeader(),
			A:   net.ParseIP(a.Data),
		}
	case dns.TypeAAAA:
		rr = &dns.AAAA{
			Hdr:  a.GetRRHeader(),
			AAAA: net.ParseIP(a.Data),
		}
	case dns.TypeCNAME:
		rr = &dns.CNAME{
			Hdr:    a.GetRRHeader(),
			Target: a.Data,
		}
	case dns.TypeNS:
		rr = &dns.NS{
			Hdr: a.GetRRHeader(),
			Ns:  a.Data,
		}
	case dns.TypeMX:
		rr = &dns.MX{
			Hdr: a.GetRRHeader(),
			Mx:  a.Data,
		}
	case dns.TypePTR:
		rr = &dns.PTR{
			Hdr: a.GetRRHeader(),
			Ptr: a.Data,
		}
	case dns.TypeSOA:
		segs := strings.Split(a.Data, " ")

		serial, _ := strconv.Atoi(segs[2])
		refresh, _ := strconv.Atoi(segs[3])
		retry, _ := strconv.Atoi(segs[4])
		expire, _ := strconv.Atoi(segs[5])
		minttl, _ := strconv.Atoi(segs[6])

		rr = &dns.SOA{
			Hdr:     a.GetRRHeader(),
			Ns:      segs[0],
			Mbox:    segs[1],
			Serial:  uint32(serial),
			Refresh: uint32(refresh),
			Retry:   uint32(retry),
			Expire:  uint32(expire),
			Minttl:  uint32(minttl),
		}
	default:
		rr = &dns.TXT{
			Hdr: dns.RR_Header{
				Name:   a.Name,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			Txt: []string{"do not support TYPE: " + dns.TypeToString[a.Type] + " currently"},
		}
	}
	return
}

// GetRRHeader 获取 rr header
// reuse code from github.com/oif/proton/gdns/response.go
func (a Answer) GetRRHeader() dns.RR_Header {
	return dns.RR_Header{
		Name:   a.Name,
		Rrtype: a.Type,
		Class:  dns.ClassINET,
		Ttl:    a.TTL,
	}
}
