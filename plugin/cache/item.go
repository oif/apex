package cache

import "github.com/miekg/dns"

type item struct {
	RecursionDesired   bool
	RecursionAvailable bool
	AuthenticatedData  bool
	CheckingDisabled   bool
	Rcode              int
	Answer             []dns.RR
	Ns                 []dns.RR
	Extra              []dns.RR

	originTTL uint32
	storedAt  uint32
}

func newItem(m *dns.Msg) *item {
	i := new(item)
	i.Rcode = m.Rcode
	i.RecursionDesired = m.RecursionDesired
	i.RecursionAvailable = m.RecursionAvailable
	i.AuthenticatedData = m.AuthenticatedData
	i.CheckingDisabled = m.CheckingDisabled

	// fmt.Printf("ans %v", m.Answer)

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
	tmp := new(dns.Msg)
	tmp.SetReply(m)
	tmp.Compress = true
	tmp.Authoritative = false
	tmp.Rcode = i.Rcode
	// m.RecursionDesired = i.RecursionDesired
	tmp.RecursionAvailable = i.RecursionAvailable
	tmp.AuthenticatedData = i.AuthenticatedData
	// m.CheckingDisabled = i.CheckingDisabled

	tmp.Answer = i.Answer
	tmp.Ns = i.Ns
	tmp.Extra = i.Extra

	tmp.CopyTo(m)
}
