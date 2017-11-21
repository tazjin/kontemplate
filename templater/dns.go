// Copyright (C) 2016-2017  Vincent Ambo <mail@tazj.in>
//
// This file is part of Kontemplate.
//
// Kontemplate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This file contains the implementation of a template function for retrieving
// IP addresses from DNS

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
