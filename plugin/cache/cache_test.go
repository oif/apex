package cache

import (
	"net"
	"testing"

	ca "github.com/oif/apex/pkg/cache"

	"github.com/miekg/dns"
)

var (
	TestQName             = "apebits.com"
	TestQType             = dns.TypeA
	TestQIP               = net.IPv4(170, 67, 54, 54)
	TestHashResult uint64 = 6082091438920753102
)

func TestCacheLRU(t *testing.T) {
	cache = ca.New(512)
	m := &dns.Msg{
		Question: []dns.Question{
			dns.Question{
				Name:  TestQName,
				Qtype: TestQType,
			},
		},
	}
	writeCache(m, TestQIP)
	if getCache(m, TestQIP) == false {
		t.Fatal("Get cache fail")
	}
}

func BenchmarkWriteCache(b *testing.B) {
	cache = ca.New(512)
	m := &dns.Msg{
		Question: []dns.Question{
			dns.Question{
				Name:  TestQName,
				Qtype: TestQType,
			},
		},
	}
	for i := 0; i < b.N; i++ {
		writeCache(m, TestQIP)
	}
}

func BenchmarkGetCache(b *testing.B) {
	cache = ca.New(512)
	m := &dns.Msg{
		Question: []dns.Question{
			dns.Question{
				Name:  TestQName,
				Qtype: TestQType,
			},
		},
	}
	writeCache(m, TestQIP)
	for i := 0; i < b.N; i++ {
		getCache(m, TestQIP)
	}
}
