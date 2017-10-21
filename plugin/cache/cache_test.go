package cache

import (
	"net"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/miekg/dns"
)

var (
	TestQName             = "apebits.com"
	TestQType             = dns.TypeA
	TestQIP               = net.IPv4(170, 67, 54, 54)
	TestHashResult uint64 = 6082091438920753102
)

func TestCacheLRU(t *testing.T) {
	var err error
	cache, err = lru.New(512)
	if err != nil {
		t.Fatal(err)
	}
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
	var err error
	cache, err = lru.New(512)
	if err != nil {
		b.Fatal(err)
	}
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
	var err error
	cache, err = lru.New(512)
	if err != nil {
		b.Fatal(err)
	}
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
