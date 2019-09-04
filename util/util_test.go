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
	"reflect"
	"testing"
)

func TestMergeWithEmptyMap(t *testing.T) {
	testMap := map[string]interface{}{
		"foo": "bar",
	}

	empty := make(map[string]interface{})

	res1 := Merge(&testMap, &empty)
	res2 := Merge(&empty, &testMap)

	if res1 != &testMap || res2 != &testMap {
		t.Error("A new map was returned incorrectly.")
		t.Fail()
	}
}

func TestMergeWithNilMap(t *testing.T) {
	testMap := map[string]interface{}{
		"foo": "bar",
	}

	res1 := Merge(&testMap, nil)
	res2 := Merge(nil, &testMap)

	if res1 != &testMap || res2 != &testMap {
		t.Error("A new map was returned incorrectly.")
		t.Fail()
	}
}

func TestMergeMaps(t *testing.T) {
	map1 := map[string]interface{}{
		"foo": "bar",
	}

	map2 := map[string]interface{}{
		"bar": "baz",
	}

	result := Merge(&map1, &map2)
	expected := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}

	if !reflect.DeepEqual(*result, expected) {
		t.Error("Maps were merged incorrectly.")
		t.Fail()
	}
}

func TestMergeMapsPrecedence(t *testing.T) {
	map1 := map[string]interface{}{
		"foo": "incorrect",
	}

	map2 := map[string]interface{}{
		"foo": "correct",
	}

	result := Merge(&map1, &map2)

	if (*result)["foo"] != "correct" {
		t.Error("Map merge precedence test failed.")
		t.Fail()
	}
}
