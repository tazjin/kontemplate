// Copyright (C) 2017  Niklas Wik <niklas.wik@nokia.com>
//
// This file is part of Kontemplate.
//
// Kontemplate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package templater

import (
	"io/ioutil"
)

//GetFromFile returns file content as string
func GetFromFile(file string) (string, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
