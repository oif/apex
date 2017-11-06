package hosts

import (
	"bufio"
	"net"
	"net/http"
	"strings"
	"time"
)

func parseHosts(url string) (map[string]net.IP, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	hosts := make(map[string]net.IP)

	lines := bufio.NewReader(resp.Body)
	for {
		line, err := lines.ReadString('\n')
		if err != nil {
			break
		}
		if len(line) < 2 || line[0] == 35 || line[0] == 32 { // if first character is '#' or space
			continue
		}
		if strings.ContainsAny(line, ":") { // ignore ipv6
			continue
		}
		//  split line
		var pair []string
		pair = strings.SplitN(line, "\t", 2)
		if pair[0] == "" || pair[1] == "" {
			continue
		}
		pair[1] = strings.TrimRight(pair[1], "\n")
		pair[1] = strings.TrimSpace(pair[1])
		hosts[pair[1]+"."] = net.ParseIP(pair[0])
	}
	return hosts, nil
}
