package mupstream

import (
	"net"
	"time"

	"github.com/miekg/dns"
)

var checkClient = &dns.Client{
	Net:          "tcp",
	ReadTimeout:  2 * time.Second,
	WriteTimeout: 2 * time.Second,
}

type upstream struct {
	addr   string
	srtt   float32
	fails  uint32
	client *dns.Client
}

func newUpstream(addr string) *upstream {
	return &upstream{
		addr: addr,
		client: &dns.Client{
			Net: "udp",
			Dialer: &net.Dialer{
				KeepAlive: time.Minute,
			},
			Timeout: 300 * time.Millisecond,
		},
	}
}

func (u *upstream) String() string {
	return u.addr
}

func (u *upstream) srttAttenuation() {
	u.srtt *= 0.98
}

func (u *upstream) forward(m *dns.Msg) (*dns.Msg, time.Duration, error) {
	resp, rtt, err := u.client.Exchange(m, u.addr)
	if err != nil { //
		u.srtt = u.srtt + 200
	} else { // success
		if rtt > 300 {
			rtt = 300
		}
		u.srtt = u.srtt*0.7 + float32(rtt)*0.3
	}
	return resp, rtt, err
}
