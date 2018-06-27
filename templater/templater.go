// Copyright (C) 2016-2017  Vincent Ambo <mail@tazj.in>
//
// This file is part of Kontemplate.
//
// Kontemplate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package templater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/tazjin/kontemplate/context"
	"github.com/tazjin/kontemplate/util"
)

const failOnMissingKeys string = "missingkey=error"

type RenderedResource struct {
	Filename string
	Rendered string
}

type RenderedResourceSet struct {
	Name      string
	Resources []RenderedResource
}

func LoadAndApplyTemplates(include *[]string, exclude *[]string, c *context.Context) ([]RenderedResourceSet, error) {
	limitedResourceSets := applyLimits(&c.ResourceSets, include, exclude)
	renderedResourceSets := make([]RenderedResourceSet, 0)

	if len(*limitedResourceSets) == 0 {
		return renderedResourceSets, fmt.Errorf("No valid resource sets included!")
	}

	for _, rs := range *limitedResourceSets {
		set, err := processResourceSet(c, &rs)

		if err != nil {
			return nil, err
		}

		renderedResourceSets = append(renderedResourceSets, *set)
	}

	return renderedResourceSets, nil
}

func processResourceSet(ctx *context.Context, rs *context.ResourceSet) (*RenderedResourceSet, error) {
	fmt.Fprintf(os.Stderr, "Loading resources for %s\n", rs.Name)

	resourcePath := path.Join(ctx.BaseDir, rs.Path)
	fileInfo, err := os.Stat(resourcePath)
	if err != nil {
		return nil, err
	}

	var files []os.FileInfo
	var resources []RenderedResource

	// Treat single-file resource paths separately from resource
	// sets containing multiple templates
	if fileInfo.IsDir() {
		// Explicitly discard this error, which will give us an empty
		// list of files instead.
		// This will end up printing a warning to the user, but it
		// won't stop the rest of the process.
		files, _ = ioutil.ReadDir(resourcePath)
		resources, err = processFiles(ctx, rs, files)
		if err != nil {
			return nil, err
		}
	} else {
		resource, err := templateFile(ctx, rs, resourcePath)
		if err != nil {
			return nil, err
		}

		resources = []RenderedResource{resource}
	}

	return &RenderedResourceSet{
		Name:      rs.Name,
		Resources: resources,
	}, nil
}

func processFiles(ctx *context.Context, rs *context.ResourceSet, files []os.FileInfo) ([]RenderedResource, error) {
	resources := make([]RenderedResource, 0)

	for _, file := range files {
		if !file.IsDir() && isResourceFile(file) {
			path := path.Join(ctx.BaseDir, rs.Path, file.Name())
			res, err := templateFile(ctx, rs, path)

			if err != nil {
				return resources, err
			}

			resources = append(resources, res)
		}
	}

	return resources, nil
}

func templateFile(ctx *context.Context, rs *context.ResourceSet, filepath string) (RenderedResource, error) {
	var resource RenderedResource

	tpl, err := template.New(path.Base(filepath)).Funcs(templateFuncs(ctx, rs)).Option(failOnMissingKeys).ParseFiles(filepath)
	if err != nil {
		return resource, fmt.Errorf("Could not load template %s: %v", filepath, err)
	}

	var b bytes.Buffer
	err = tpl.Execute(&b, rs.Values)
	if err != nil {
		return resource, fmt.Errorf("Error while templating %s: %v", filepath, err)
	}

	resource = RenderedResource{
		Filename: path.Base(filepath),
		Rendered: b.String(),
	}

	return resource, nil
}

// Applies the limits of explicitly included or excluded resources and returns the updated resource set.
// Exclude takes priority over include
func applyLimits(rs *[]context.ResourceSet, include *[]string, exclude *[]string) *[]context.ResourceSet {
	if len(*include) == 0 && len(*exclude) == 0 {
		return rs
	}

	// Exclude excluded resource sets
	excluded := make([]context.ResourceSet, 0)
	for _, r := range *rs {
		if !matchesResourceSet(exclude, &r) {
			excluded = append(excluded, r)
		}
	}

	// Include included resource sets
	if len(*include) == 0 {
		return &excluded
	}
	included := make([]context.ResourceSet, 0)
	for _, r := range excluded {
		if matchesResourceSet(include, &r) {
			included = append(included, r)
		}
	}

	return &included
}

// Check whether an include/exclude string slice matches a resource set
func matchesResourceSet(s *[]string, rs *context.ResourceSet) bool {
	for _, r := range *s {
		r = strings.TrimSuffix(r, "/")
		if r == rs.Name || r == rs.Parent {
			return true
		}
	}

	return false
}

func templateFuncs(c *context.Context, rs *context.ResourceSet) template.FuncMap {
	m := sprig.TxtFuncMap()
	m["json"] = func(data interface{}) string {
		b, _ := json.Marshal(data)
		return string(b)
	}
	m["passLookup"] = GetFromPass
	m["gitHEAD"] = func() (string, error) {
		out, err := exec.Command("git", "-C", c.BaseDir, "rev-parse", "HEAD").Output()
		if err != nil {
			return "", err
		}
		output := strings.TrimSpace(string(out))
		return output, nil
	}
	m["lookupIPAddr"] = GetIPsFromDNS
	m["insertFile"] = func(file string) (string, error) {
		data, err := ioutil.ReadFile(path.Join(rs.Path, file))
		if err != nil {
			return "", err
		}

		return string(data), nil
	}
	m["default"] = func(defaultVal interface{}, varName string) interface{} {
		if val, ok := rs.Values[varName]; ok {
			return val
		}

		return defaultVal
	}
	return m
}

// Checks whether a file is a resource file (i.e. is YAML or JSON) and not a default values file.
func isResourceFile(f os.FileInfo) bool {
	for _, defaultFile := range util.DefaultFilenames {
		if f.Name() == defaultFile {
			return false
		}
	}

	return strings.HasSuffix(f.Name(), "yaml") ||
		strings.HasSuffix(f.Name(), "yml") ||
		strings.HasSuffix(f.Name(), "json")
}
