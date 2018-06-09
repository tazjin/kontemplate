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
				Path: "some-api",
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
		BaseDir: "testdata",
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
				Path: "some-api",
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
				Path: "collection/nested",
				Values: map[string]interface{}{
					"lizards": "good",
					"globalVar": "lizards",
				},
				Include: nil,
				Parent:  "collection",
			},
		},
		BaseDir: "testdata",
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
				Path: "parent/child",
				Values: map[string]interface{}{
					"foo": "bar",
					"bar": "baz",
				},
				Include: nil,
				Parent:  "parent",
			},
		},
		BaseDir: "testdata",
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
				Path: "parent/child",
				Values: map[string]interface{}{
					"foo": "newvalue",
				},
				Include: nil,
				Parent:  "parent",
			},
		},
		BaseDir: "testdata",
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

func TestValuesOverride(t *testing.T) {
	ctx, err := LoadContext("testdata/import-vars-override.yaml", &noExplicitVars)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expected := map[string]interface{}{
		"override": float64(3),
		"music": map[string]interface{}{
			"artist": "Pallida",
			"track":  "Tractor Beam",
		},
		"place": "Oslo",
		"globalVar": "very global!",
	}

	if !reflect.DeepEqual(ctx.ResourceSets[0].Values, expected) {
		t.Error("Expected overrides after loading imports did not match!")
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
				Path: "some-api",
				Values: map[string]interface{}{
					"location": "europe",
				},
				Include: nil,
				Parent:  "",
			},
			{
				Name: "some-api-asia",
				Path: "some-api",
				Values: map[string]interface{}{
					"location": "asia",
				},
				Include: nil,
				Parent:  "",
			},
		},
		BaseDir: "testdata",
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
				Path:   "parent-path/child-path",
				Parent: "parent",
				Values: make(map[string]interface{}, 0),
			},
		},
		BaseDir: "testdata",
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
