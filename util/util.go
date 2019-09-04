// Copyright (C) 2016-2019  Vincent Ambo <mail@tazj.in>
//
// This file is part of Kontemplate.
//
// Kontemplate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package util

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Filenames excluded from templating for the purpose of containing default variable values inside a resource set.
var DefaultFilenames []string = []string{"default.yml", "default.yaml", "default.json"}

// Merges two maps together. Values from the second map override values in the first map.
// The returned map is new if anything was changed.
func Merge(in1 *map[string]interface{}, in2 *map[string]interface{}) *map[string]interface{} {
	if in1 == nil || len(*in1) == 0 {
		return in2
	}

	if in2 == nil || len(*in2) == 0 {
		return in1
	}

	new := make(map[string]interface{})
	for k, v := range *in1 {
		new[k] = v
	}

	for k, v := range *in2 {
		new[k] = v
	}

	return &new
}

// Loads either a YAML or JSON file from the specified path and
// deserialises it into the provided interface.
func LoadData(filename string, addr interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, addr)
	if err != nil {
		return err
	}

	return nil
}
