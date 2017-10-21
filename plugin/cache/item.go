package cache

import "github.com/miekg/dns"

type item struct {
	Rcode              int
	Authoritative      bool
	AuthenticatedData  bool
	RecursionAvailable bool
	Answer             []dns.RR
	Ns                 []dns.RR
	Extra              []dns.RR
}

func newItem(m *dns.Msg) *item {
	i := new(item)
	i.Rcode = m.Rcode
	i.Authoritative = m.Authoritative
	i.AuthenticatedData = m.AuthenticatedData
	i.RecursionAvailable = m.RecursionAvailable
	i.Answer = m.Answer
	i.Ns = m.Ns
	i.Extra = make([]dns.RR, len(m.Extra))
	// Don't copy OPT record as these are hop-by-hop.
	j := 0
	for _, e := range m.Extra {
		if e.Header().Rrtype == dns.TypeOPT {
			continue
		}
		i.Extra[j] = e
		j++
	}
	i.Extra = i.Extra[:j]

	return i
}

func (i *item) replyToMsg(m *dns.Msg) {
	m.Authoritative = false
	m.AuthenticatedData = i.AuthenticatedData
	m.RecursionAvailable = i.RecursionAvailable
	m.Rcode = i.Rcode
	m.Compress = true

	m.Answer = i.Answer
	m.Ns = i.Ns
	m.Extra = i.Extra

}
