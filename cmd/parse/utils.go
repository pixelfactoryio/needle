package parse

import (
	"context"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
)

func getHostFile(file string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, file, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return bs, nil
}

func WriteHosts(hosts []string, file string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, host := range hosts {
		if _, err := f.WriteString(host + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func filterHosts(field string) bool {
	// filter unwanted hosts
	return slices.Contains([]string{
		"broadcasthost",
		"ip6-allhosts",
		"ip6-allnodes",
		"ip6-allrouters",
		"ip6-localhost",
		"ip6-localnet",
		"ip6-loopback",
		"ip6-mcastprefix",
		"local",
		"localhost",
		"localhost.localdomain",
		"localhost4",
		"localhost4.localdomain4",
		"localhost6",
		"localhost6.localdomain6",
		"localdomain",
		"localdomain.local",
		"localdomain4",
		"0.0.0.0",
	}, field)
}

func Hosts(bs []byte, ipv4, ipv6 string) []string {
	hosts := []string{}

	lines := strings.Split(string(bs), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		if len(fields) > 2 {
			continue
		}

		if len(fields) == 1 {
			if filterHosts(fields[0]) {
				continue
			}

			hosts = append(hosts, ipv4+" "+fields[0], ipv6+" "+fields[0])
			continue
		}

		if len(fields) == 2 {
			if filterHosts(fields[1]) {
				continue
			}

			hosts = append(hosts, ipv4+" "+fields[1], ipv6+" "+fields[1])
			continue
		}
	}

	return hosts
}
