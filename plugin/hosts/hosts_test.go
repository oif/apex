package hosts

import "testing"

func TestLoadHosts(t *testing.T) {
	_, err := parseHosts("https://raw.githubusercontent.com/racaljk/hosts/master/hosts")
	if err != nil {
		t.Fatal(err)
	}
}
