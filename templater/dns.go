// Copyright (C) 2016-2019  Vincent Ambo <mail@tazj.in>
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
	"net"
	"os"
)

func GetIPsFromDNS(host string) ([]interface{}, error) {
	fmt.Fprintf(os.Stderr, "Attempting to look up IP for %s in DNS\n", host)
	ips, err := net.LookupIP(host)

	if err != nil {
		return nil, fmt.Errorf("IP address lookup failed: %v", err)
	}

	var result []interface{} = make([]interface{}, len(ips))
	for i, ip := range ips {
		result[i] = ip
	}

	return result, nil
}
