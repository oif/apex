package gdns

import "encoding/json"

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
	return r.Status == 0, r.Comment
}
