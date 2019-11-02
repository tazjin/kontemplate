// Copyright (C) 2016-2019  Vincent Ambo <mail@tazj.in>
//
// This file is part of Kontemplate.
//
// Kontemplate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package templater

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tazjin/kontemplate/context"
)

func TestApplyNoLimits(t *testing.T) {
	resources := []context.ResourceSet{
		{
			Name: "testResourceSet1",
		},
		{
			Name: "testResourceSet2",
		},
	}

	result := applyLimits(&resources, &[]string{}, &[]string{})

	if !reflect.DeepEqual(resources, *result) {
		t.Error("Resource set slice changed, but shouldn't have.")
		t.Errorf("Expected: %v\nResult: %v\n", resources, *result)
		t.Fail()
	}
}

func TestApplyIncludeLimits(t *testing.T) {
	resources := []context.ResourceSet{
		{
			Name: "testResourceSet1",
		},
		{
			Name: "testResourceSet2",
		},
		{
			Name:   "testResourceSet3",
			Parent: "included",
		},
	}

	includes := []string{"testResourceSet1", "included"}

	result := applyLimits(&resources, &includes, &[]string{})

	expected := []context.ResourceSet{
		{
			Name: "testResourceSet1",
		},
		{
			Name:   "testResourceSet3",
			Parent: "included",
		},
	}

	if !reflect.DeepEqual(expected, *result) {
		t.Error("Result does not contain expected resource sets.")
		t.Errorf("Expected: %v\nResult: %v\n", expected, *result)
		t.Fail()
	}
}

func TestApplyExcludeLimits(t *testing.T) {
	resources := []context.ResourceSet{
		{
			Name: "testResourceSet1",
		},
		{
			Name: "testResourceSet2",
		},
		{
			Name:   "testResourceSet3",
			Parent: "included",
		},
	}

	exclude := []string{"testResourceSet2"}

	result := applyLimits(&resources, &[]string{}, &exclude)

	expected := []context.ResourceSet{
		{
			Name: "testResourceSet1",
		},
		{
			Name:   "testResourceSet3",
			Parent: "included",
		},
	}

	if !reflect.DeepEqual(expected, *result) {
		t.Error("Result does not contain expected resource sets.")
		t.Errorf("Expected: %v\nResult: %v\n", expected, *result)
		t.Fail()
	}
}

func TestApplyLimitsExcludeIncludePrecedence(t *testing.T) {
	resources := []context.ResourceSet{
		{
			Name:   "collection/nested1",
			Parent: "collection",
		},
		{
			Name:   "collection/nested2",
			Parent: "collection",
		},
		{
			Name:   "collection/nested3",
			Parent: "collection",
		},
		{
			Name: "something-else",
		},
	}

	include := []string{"collection"}
	exclude := []string{"collection/nested2"}

	result := applyLimits(&resources, &include, &exclude)

	expected := []context.ResourceSet{
		{
			Name:   "collection/nested1",
			Parent: "collection",
		},
		{
			Name:   "collection/nested3",
			Parent: "collection",
		},
	}

	if !reflect.DeepEqual(expected, *result) {
		t.Error("Result does not contain expected resource sets.")
		t.Errorf("Expected: %v\nResult: %v\n", expected, *result)
		t.Fail()
	}
}

func TestFailOnMissingKeys(t *testing.T) {
	ctx := context.Context{}
	resourceSet := context.ResourceSet{}

	_, err := templateFile(&ctx, &resourceSet, "testdata/test-template.txt")

	if err == nil {
		t.Errorf("Template with missing keys should have failed.\n")
		t.Fail()
	}

	if !strings.Contains(err.Error(), "map has no entry for key \"testName\"") {
		t.Errorf("Templating failed with unexpected error: %v\n", err)
	}
}

func TestDefaultTemplateFunction(t *testing.T) {
	ctx := context.Context{}
	resourceSet := context.ResourceSet{}

	res, err := templateFile(&ctx, &resourceSet, "testdata/test-default.txt")

	if err != nil {
		t.Errorf("Templating with default values should have succeeded.\n")
		t.Fail()
	}

	if res.Rendered != "defaultValue\n" {
		t.Error("Result does not contain expected rendered default value.")
		t.Fail()
	}
}

func TestInsertTemplateFunction(t *testing.T) {
	ctx := context.Context{}
	resourceSet := context.ResourceSet{
		Path: "testdata",
		Values: map[string]interface{}{
			"testName": "TestInsertTemplateFunction",
		},
	}

	res, err := templateFile(&ctx, &resourceSet, "testdata/test-insertTemplate.txt")

	if err != nil {
		t.Error(err)
		t.Errorf("Templating with an insertTemplate call should have succeeded.\n")
		t.Fail()
	}
	if res.Rendered != "Inserting \"Template for test TestInsertTemplateFunction\".\n" {
		t.Error("Result does not contain expected rendered template value.")
		t.Error(res.Rendered)
		t.Fail()
	}
}

func TestInsertTemplateFunctionWithArgs(t *testing.T) {
	ctx := context.Context{}
	resourceSet := context.ResourceSet{
		Path: "testdata",
		Values: map[string]interface{}{
			"testName": "TestInsertTemplateFunctionWithArgs",
			"foo":      "bar",
		},
	}

	res, err := templateFile(&ctx, &resourceSet, "testdata/test-insertTemplateWithArgs.txt")

	if err != nil {
		t.Error(err)
		t.Errorf("Templating with an insertTemplate call should have succeeded.\n")
		t.Fail()
	}

	expected := "Overriding \"Template for test override\".\nOriginal \"Template for test TestInsertTemplateFunctionWithArgs\"."
	if res.Rendered != expected {
		t.Error("Result does not contain expected rendered template value.")
		t.Error(res.Rendered)

		t.Fail()
	}
}
