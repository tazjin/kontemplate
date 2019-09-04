// Copyright (C) 2016-2017  Vincent Ambo <mail@tazj.in>
//
// This file is part of Kontemplate.
//
// Kontemplate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package context

import (
	"reflect"
	"testing"
)

var noExplicitVars []string = make([]string, 0)

func TestLoadFlatContextFromFile(t *testing.T) {
	ctx, err := LoadContext("testdata/flat-test.yaml", &noExplicitVars)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		Global: map[string]interface{}{
			"globalVar": "lizards",
		},
		ResourceSets: []ResourceSet{
			{
				Name: "some-api",
				Path: "testdata/some-api",
				Values: map[string]interface{}{
					"apiPort":          float64(4567), // yep!
					"importantFeature": true,
					"version":          "1.0-0e6884d",
					"globalVar":        "lizards",
				},
				Include: nil,
				Parent:  "",
			},
		},
		BaseDir:      "testdata",
		ImportedVars: make(map[string]interface{}, 0),
		ExplicitVars: make(map[string]interface{}, 0),
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded context and expected context did not match")
		t.Fail()
	}
}

func TestLoadContextWithArgs(t *testing.T) {
	ctx, err := LoadContext("testdata/flat-with-args-test.yaml", &noExplicitVars)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name:   "some-api",
				Path:   "testdata/some-api",
				Values: make(map[string]interface{}, 0),
				Args: []string{
					"--as=some-user",
					"--as-group=hello:world",
					"--as-banana",
					"true",
				},
				Include: nil,
				Parent:  "",
			},
		},
		BaseDir:      "testdata",
		ImportedVars: make(map[string]interface{}, 0),
		ExplicitVars: make(map[string]interface{}, 0),
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded context and expected context did not match")
		t.Fail()
	}
}

func TestLoadContextWithResourceSetCollections(t *testing.T) {
	ctx, err := LoadContext("testdata/collections-test.yaml", &noExplicitVars)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		Global: map[string]interface{}{
			"globalVar": "lizards",
		},
		ResourceSets: []ResourceSet{
			{
				Name: "some-api",
				Path: "testdata/some-api",
				Values: map[string]interface{}{
					"apiPort":          float64(4567), // yep!
					"importantFeature": true,
					"version":          "1.0-0e6884d",
					"globalVar":        "lizards",
				},
				Include: nil,
				Parent:  "",
			},
			{
				Name: "collection/nested",
				Path: "testdata/collection/nested",
				Values: map[string]interface{}{
					"lizards":   "good",
					"globalVar": "lizards",
				},
				Include: nil,
				Parent:  "collection",
			},
		},
		BaseDir:      "testdata",
		ImportedVars: make(map[string]interface{}, 0),
		ExplicitVars: make(map[string]interface{}, 0),
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded context and expected context did not match")
		t.Fail()
	}

}

func TestSubresourceVariableInheritance(t *testing.T) {
	ctx, err := LoadContext("testdata/parent-variables.yaml", &noExplicitVars)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name: "parent/child",
				Path: "testdata/parent/child",
				Values: map[string]interface{}{
					"foo": "bar",
					"bar": "baz",
				},
				Include: nil,
				Parent:  "parent",
			},
		},
		BaseDir:      "testdata",
		ImportedVars: make(map[string]interface{}, 0),
		ExplicitVars: make(map[string]interface{}, 0),
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded and expected context did not match")
		t.Fail()
	}
}

func TestSubresourceVariableInheritanceOverride(t *testing.T) {
	ctx, err := LoadContext("testdata/parent-variable-override.yaml", &noExplicitVars)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name: "parent/child",
				Path: "testdata/parent/child",
				Values: map[string]interface{}{
					"foo": "newvalue",
				},
				Include: nil,
				Parent:  "parent",
			},
		},
		BaseDir:      "testdata",
		ImportedVars: make(map[string]interface{}, 0),
		ExplicitVars: make(map[string]interface{}, 0),
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded and expected context did not match")
		t.Fail()
	}
}

func TestDefaultValuesLoading(t *testing.T) {
	ctx, err := LoadContext("testdata/default-loading.yaml", &noExplicitVars)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	rs := ctx.ResourceSets[0]
	if rs.Values["defaultValues"] != "loaded" {
		t.Errorf("Default values not loaded from YAML file")
		t.Fail()
	}

	if rs.Values["override"] != "notAtAll" {
		t.Error("Default values should not override other values")
		t.Fail()
	}
}

func TestImportValuesLoading(t *testing.T) {
	ctx, err := LoadContext("testdata/import-vars-simple.yaml", &noExplicitVars)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := map[string]interface{}{
		"override": "true",
		"music": map[string]interface{}{
			"artist": "Pallida",
			"track":  "Tractor Beam",
		},
	}

	if !reflect.DeepEqual(ctx.ImportedVars, expected) {
		t.Error("Expected imported values after loading imports did not match!")
		t.Fail()
	}
}

func TestExplicitPathLoading(t *testing.T) {
	ctx, err := LoadContext("testdata/explicit-path.yaml", &noExplicitVars)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name: "some-api-europe",
				Path: "testdata/some-api",
				Values: map[string]interface{}{
					"location": "europe",
				},
				Include: nil,
				Parent:  "",
			},
			{
				Name: "some-api-asia",
				Path: "testdata/some-api",
				Values: map[string]interface{}{
					"location": "asia",
				},
				Include: nil,
				Parent:  "",
			},
		},
		BaseDir:      "testdata",
		ImportedVars: make(map[string]interface{}, 0),
		ExplicitVars: make(map[string]interface{}, 0),
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded context and expected context did not match")
		t.Fail()
	}
}

func TestExplicitSubresourcePathLoading(t *testing.T) {
	ctx, err := LoadContext("testdata/explicit-subresource-path.yaml", &noExplicitVars)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := Context{
		Name: "k8s.prod.mydomain.com",
		ResourceSets: []ResourceSet{
			{
				Name:   "parent/child",
				Path:   "testdata/parent-path/child-path",
				Parent: "parent",
				Values: make(map[string]interface{}, 0),
			},
		},
		BaseDir:      "testdata",
		ImportedVars: make(map[string]interface{}, 0),
		ExplicitVars: make(map[string]interface{}, 0),
	}

	if !reflect.DeepEqual(*ctx, expected) {
		t.Error("Loaded context and expected context did not match")
		t.Fail()
	}
}

func TestSetVariablesFromArguments(t *testing.T) {
	vars := []string{"version=some-service-version"}
	ctx, _ := LoadContext("testdata/default-loading.yaml", &vars)

	if version := ctx.ExplicitVars["version"]; version != "some-service-version" {
		t.Errorf(`Expected variable "version" to have value "some-service-version" but was "%s"`, version)
	}
}

func TestSetInvalidVariablesFromArguments(t *testing.T) {
	vars := []string{"version: some-service-version"}
	_, err := LoadContext("testdata/default-loading.yaml", &vars)

	if err == nil {
		t.Error("Expected invalid variable to return an error")
	}
}

// This test ensures that variables are merged in the correct order.
// Please consult the test data in `testdata/merging`.
func TestValueMergePrecedence(t *testing.T) {
	cliVars:= []string{"cliVar=cliVar"}
	ctx, _ := LoadContext("testdata/merging/context.yaml", &cliVars)

	expected := map[string]interface{}{
		"defaultVar": "defaultVar",
		"importVar": "importVar",
		"globalVar": "globalVar",
		"includeVar": "includeVar",
		"cliVar": "cliVar",
	}

	result := ctx.ResourceSets[0].Values

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Merged values did not match expected result: \n%v", result)
		t.Fail()
	}
}
