package cache

import (
	"encoding/binary"
	"hash/fnv"
	"net"

	lru "github.com/hashicorp/golang-lru"
	"github.com/miekg/dns"
)

var cache *lru.Cache // 解析缓存

func key(m *dns.Msg, clientIP net.IP) uint64 {
	if m.Truncated {
		return 0
	}
	return hash(m.Question[0].Name, m.Question[0].Qtype, clientIP)
}

func hash(qname string, qtype uint16, qip []byte) uint64 {
	h := fnv.New64()
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, qtype)
	h.Write(b)
	var c byte
	for i := range qname {
		c = qname[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		h.Write([]byte{c})
	}
	h.Write(qip)
	return h.Sum64()
}

func writeCache(m *dns.Msg, ip net.IP) {
	if len(m.Question) > 0 {
		if key := key(m, ip); key != 0 {
			cache.Add(key, newItem(m)) // if write failed, just ignore it
		}
	}
}

// return true if hit cache
func getCache(m *dns.Msg, ip net.IP) bool {
	if key := key(m, ip); key != 0 {
		cached, ok := cache.Get(key)
		if ok {
			cached.(*item).replyToMsg(m)
			return true
		}
	}
	// miss cache
	return false
}
