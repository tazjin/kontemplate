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
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"

	"github.com/tazjin/kontemplate/util"
)

type ResourceSet struct {
	// Name of the resource set. This can be used in include/exclude statements during kontemplate runs.
	Name string `json:"name"`

	// Path to the folder containing the files for this resource set. This defaults to the value of the 'name' field
	// if unset.
	Path string `json:"path"`

	// Values to include when interpolating resources from this resource set.
	Values map[string]interface{} `json:"values"`

	// Nested resource sets to include
	Include []ResourceSet `json:"include"`

	// Parent resource set for flattened resource sets. Should not be manually specified.
	Parent string
}

type Context struct {
	// The name of the kubectl context
	Name string `json:"context"`

	// Global variables that should be accessible by all resource sets
	Global map[string]interface{} `json:"global"`

	// File names of YAML or JSON files including extra variables that should be globally accessible
	VariableImports []string `json:"import"`

	// The resource sets to include in this context
	ResourceSets []ResourceSet `json:"include"`

	// This field represents the absolute path to the context base directory and should not be manually specified.
	BaseDir string
}

func contextLoadingError(filename string, cause error) error {
	return fmt.Errorf("Context loading failed on file %s due to: \n%v", filename, cause)
}

// Attempt to load and deserialise a Context from the specified file.
func LoadContextFromFile(filename string) (*Context, error) {
	var c Context
	err := util.LoadJsonOrYaml(filename, &c)

	if err != nil {
		return nil, contextLoadingError(filename, err)
	}

	c.ResourceSets = flattenPrepareResourceSetPaths(&c.ResourceSets)
	c.BaseDir = path.Dir(filename)
	c.ResourceSets = loadAllDefaultValues(&c)

	err = c.loadImportedVariables()
	if err != nil {
		return nil, contextLoadingError(filename, err)
	}

	err = c.loadStdinVariables()
	if err != nil {
		return nil, contextLoadingError("stdin", err)
	}

	return &c, nil
}

// Kontemplate supports specifying additional variable files with the `import` keyword. This function loads those
// variable files and merges them together with the context's other global variables.
func (ctx *Context) loadImportedVariables() error {
	for _, file := range ctx.VariableImports {
		var importedVars map[string]interface{}
		err := util.LoadJsonOrYaml(path.Join(ctx.BaseDir, file), &importedVars)

		if err != nil {
			return err
		}

		ctx.Global = *util.Merge(&ctx.Global, &importedVars)
	}

	return nil
}

// As a bonus load additional variables from data piped to stdin as if they had been added with `import`.
func (ctx *Context) loadStdinVariables() error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// Stdin is a TTY, skip this
		return nil
	}

	file, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var stdinVars map[string]interface{}
	err = yaml.Unmarshal(file, &stdinVars)
	if err != nil {
		return err
	}
	ctx.Global = *util.Merge(&ctx.Global, &stdinVars)
	return err
}

// Correctly prepares the file paths for resource sets by inferring implicit paths and flattening resource set
// collections, i.e. resource sets that themselves have an additional 'include' field set.
// Those will be regarded as a short-hand for including multiple resource sets from a subfolder.
// See https://github.com/tazjin/kontemplate/issues/9 for more information.
func flattenPrepareResourceSetPaths(rs *[]ResourceSet) []ResourceSet {
	flattened := make([]ResourceSet, 0)

	for _, r := range *rs {
		// If a path is not explicitly specified it should default to the resource set name.
		// This is also the classic behaviour prior to kontemplate 1.2
		if r.Path == "" {
			r.Path = r.Name
		}

		if len(r.Include) == 0 {
			flattened = append(flattened, r)
		} else {
			for _, subResourceSet := range r.Include {
				if subResourceSet.Path == "" {
					subResourceSet.Path = subResourceSet.Name
				}

				subResourceSet.Parent = r.Name
				subResourceSet.Name = path.Join(r.Name, subResourceSet.Name)
				subResourceSet.Path = path.Join(r.Path, subResourceSet.Path)
				subResourceSet.Values = *util.Merge(&r.Values, &subResourceSet.Values)
				flattened = append(flattened, subResourceSet)
			}
		}
	}

	return flattened
}

func loadAllDefaultValues(c *Context) []ResourceSet {
	updated := make([]ResourceSet, len(c.ResourceSets))

	for i, rs := range c.ResourceSets {
		merged := loadDefaultValues(&rs, c)
		rs.Values = *merged
		updated[i] = rs
	}

	return updated
}

// Loads and merges default values for a resource set collection from path/to/set/default.{json|yaml}.
// YAML takes precedence over JSON.
// Default values in resource set collections have the lowest priority possible.
func loadDefaultValues(rs *ResourceSet, c *Context) *map[string]interface{} {
	var defaultVars map[string]interface{}

	for _, filename := range util.DefaultFilenames {
		err := util.LoadJsonOrYaml(path.Join(c.BaseDir, rs.Path, filename), &defaultVars)
		if err == nil {
			return util.Merge(&defaultVars, &rs.Values)
		}
	}

	// The actual error is not inspected here. The reasoning for this is that in case of serious problems (e.g.
	// permission issues with the folder / folder not existing) failure will occur a bit later anyways.
	// Otherwise we'd have to differentiate between file-not-found-errors (no default values specified) and other
	// errors here.
	return &rs.Values
}
