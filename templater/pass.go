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
// variables from 'pass', the standard UNIX password manager.

package templater

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetFromPass(key string) (string, error) {
	fmt.Fprintf(os.Stderr, "Attempting to look up %s in pass\n", key)
	pass := exec.Command("pass", "show", key)

	output, err := pass.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Pass lookup failed: %s (%v)", output, err)
	}

	trimmed := strings.TrimSpace(string(output))

	return trimmed, nil
}
