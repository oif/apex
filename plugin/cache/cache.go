package cache

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coocood/freecache"
	"github.com/miekg/dns"
)

// CacheKeyFormat the format of cache store key
const CacheKeyFormat = "%s_%s_%s" // name, type, subnet
var cache *freecache.Cache        // 解析缓存

// 根据格式生成 key
func getKey(name, qtype, ip string) string {
	dot := strings.LastIndex(ip, ".")
	if dot > 0 { // v4
		return fmt.Sprintf(CacheKeyFormat, name, qtype, ip[0:dot])
	}

	colon := strings.LastIndex(ip, ":")
	if colon > 0 { // v6
		return fmt.Sprintf(CacheKeyFormat, name, qtype, ip[0:colon])
	}
	return ""
}

// 增加 DNS 缓存
func addDNSCache(m *dns.Msg, ip string) error {
	if len(m.Answer) > 0 {
		packed, err := m.Pack()
		if err != nil {
			return err
		}
		key := getKey(m.Question[0].Name, dns.TypeToString[m.Question[0].Qtype], ip)
		if key == "" {
			return errors.New("invalid ip " + ip)
		}
		cache.Set([]byte(key), packed, int(m.Answer[0].Header().Ttl))
	}
	return nil
}

// 获取 DNS 缓存
func getDNSCache(m *dns.Msg, ip string) error {
	temp := &dns.Msg{}
	key := getKey(m.Question[0].Name, dns.TypeToString[m.Question[0].Qtype], ip)
	if key == "" {
		return errors.New("invalid ip" + ip)
	}
	got, err := cache.Get([]byte(key))
	if err != nil {
		return err
	}
	// 有缓存
	err = temp.Unpack(got)
	if err != nil {
		return err
	}

	m.Ns = temp.Ns
	m.Answer = temp.Answer
	m.Extra = temp.Extra

	// 转换成功
	return nil
}
