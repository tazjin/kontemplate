// This file contains the implementation of a template function for retrieving IP addresses from DNS
package templater

import (
	"fmt"
	"github.com/polydawn/meep"
	"net"
	"os"
)

type DNSError struct {
	meep.TraitAutodescribing
	meep.TraitCausable
	Output string
}

func GetIPsFromDNS(host string) ([]interface{}, error) {
	fmt.Fprintf(os.Stderr, "Attempting to look up IP for %s in DNS\n", host)
	ips, err := net.LookupIP(host)

	if err != nil {
		return nil, meep.New(
			&DNSError{Output: "IP address lookup failed"},
			meep.Cause(err),
		)
	}

	var result []interface{} = make([]interface{}, len(ips))
	for i, ip := range ips {
		result[i] = ip
	}

	return result, nil
}
