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
	"path"
	"strings"

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

	// Args to pass on to kubectl for this resource set.
	Args []string `json:"args"`

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
	VariableImportFiles []string `json:"import"`

	// The resource sets to include in this context
	ResourceSets []ResourceSet `json:"include"`

	// Variables imported from additional files
	ImportedVars map[string]interface{}

	// Explicitly set variables (via `--var`) that should override all others
	ExplicitVars map[string]interface{}

	// This field represents the absolute path to the context base directory and should not be manually specified.
	BaseDir string
}

func contextLoadingError(filename string, cause error) error {
	return fmt.Errorf("Context loading failed on file %s due to: \n%v", filename, cause)
}

// Attempt to load and deserialise a Context from the specified file.
func LoadContext(filename string, explicitVars *[]string) (*Context, error) {
	var ctx Context
	err := util.LoadData(filename, &ctx)

	if err != nil {
		return nil, contextLoadingError(filename, err)
	}

	ctx.BaseDir = path.Dir(filename)

	// Prepare the resource sets by resolving parents etc.
	ctx.ResourceSets = flattenPrepareResourceSetPaths(&ctx.BaseDir, &ctx.ResourceSets)

	// Add variables explicitly specified on the command line
	ctx.ExplicitVars, err = loadExplicitVars(explicitVars)
	if err != nil {
		return nil, fmt.Errorf("Error setting explicit variables: %v\n", err)
	}

	// Add variables loaded from import files
	ctx.ImportedVars, err = ctx.loadImportedVariables()
	if err != nil {
		return nil, contextLoadingError(filename, err)
	}

	// Merge variables defined at different levels. The
	// `mergeContextValues` function is documented with the merge
	// hierarchy.
	ctx.ResourceSets = ctx.mergeContextValues()

	if err != nil {
		return nil, contextLoadingError(filename, err)
	}

	return &ctx, nil
}

// Kontemplate supports specifying additional variable files with the
// `import` keyword. This function loads those variable files and
// merges them together with the context's other global variables.
func (ctx *Context) loadImportedVariables() (map[string]interface{}, error) {
	allImportedVars := make(map[string]interface{})

	for _, file := range ctx.VariableImportFiles {
		// Ensure that the filename is not merged with the baseDir if
		// it is set to an absolute path.
		var filePath string
		if path.IsAbs(file) {
			filePath = file
		} else {
			filePath = path.Join(ctx.BaseDir, file)
		}

		var importedVars map[string]interface{}
		err := util.LoadData(filePath, &importedVars)

		if err != nil {
			return nil, err
		}

		allImportedVars = *util.Merge(&allImportedVars, &importedVars)
	}

	return allImportedVars, nil
}

// Correctly prepares the file paths for resource sets by inferring implicit paths and flattening resource set
// collections, i.e. resource sets that themselves have an additional 'include' field set.
// Those will be regarded as a short-hand for including multiple resource sets from a subfolder.
// See https://github.com/tazjin/kontemplate/issues/9 for more information.
func flattenPrepareResourceSetPaths(baseDir *string, rs *[]ResourceSet) []ResourceSet {
	flattened := make([]ResourceSet, 0)

	for _, r := range *rs {
		// If a path is not explicitly specified it should default to the resource set name.
		// This is also the classic behaviour prior to kontemplate 1.2
		if r.Path == "" {
			r.Path = r.Name
		}

		// Paths are made absolute by resolving them relative to the context base,
		// unless absolute paths were specified.
		if !path.IsAbs(r.Path) {
			r.Path = path.Join(*baseDir, r.Path)
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

// Merges the context and resource set variables according in the
// desired precedence order.
//
// For now the reasoning behind the merge order is from least specific
// in relation to the cluster configuration, which means that the
// precedence is (in ascending order):
//
// 1. Default values in resource sets.
// 2. Values imported from files (via `import:`)
// 3. Global values in a cluster configuration
// 4. Values set in a resource set's `include`-section
// 5. Explicit values set on the CLI (`--var`)
//
// For a discussion on the reasoning behind this order, please consult
// https://github.com/tazjin/kontemplate/issues/142
func (ctx *Context) mergeContextValues() []ResourceSet {
	updated := make([]ResourceSet, len(ctx.ResourceSets))

	// Merging has to happen separately for every individual
	// resource set to make use of the default values:
	for i, rs := range ctx.ResourceSets {
		// Begin by loading default values from the resource
		// sets configuration.
		//
		// Resource sets are used across different cluster
		// contexts and the default values in them have the
		// lowest precedence.
		defaultValues := loadDefaultValues(&rs, ctx)

		// Continue by merging default values with values
		// imported from external files. Those values are also
		// used across cluster contexts, but have higher
		// precedence than defaults.
		merged := util.Merge(defaultValues, &ctx.ImportedVars)

		// Merge global values defined in the cluster context:
		merged = util.Merge(merged, &ctx.Global)

		// Merge values configured in the resource set's
		// `include` section:
		merged = util.Merge(merged, &rs.Values)

		// Merge values defined explicitly on the CLI:
		merged = util.Merge(merged, &ctx.ExplicitVars)

		// Continue with the newly merged resource set:
		rs.Values = *merged
		updated[i] = rs
	}

	return updated
}

// Loads default values for a resource set collection from
// path/to/set/default.{json|yaml}.
func loadDefaultValues(rs *ResourceSet, c *Context) *map[string]interface{} {
	var defaultVars map[string]interface{}

	for _, filename := range util.DefaultFilenames {
		err := util.LoadData(path.Join(rs.Path, filename), &defaultVars)
		if err == nil {
			return &defaultVars
		}
	}

	// The actual error is not inspected here. The reasoning for
	// this is that in case of serious problems (e.g. permission
	// issues with the folder / folder not existing) failure will
	// occur a bit later anyways.
	//
	// Otherwise we'd have to differentiate between
	// file-not-found-errors (no default values specified) and
	// other errors here.
	return &rs.Values
}

// Prepares the variables specified explicitly via `--var` when
// executing kontemplate for adding to the context.
func loadExplicitVars(vars *[]string) (map[string]interface{}, error) {
	explicitVars := make(map[string]interface{}, len(*vars))

	for _, v := range *vars {
		varParts := strings.SplitN(v, "=", 2)
		if len(varParts) != 2 {
			return nil, fmt.Errorf(`invalid explicit variable provided (%s), name and value should be separated with "="`, v)
		}

		explicitVars[varParts[0]] = varParts[1]
	}

	return explicitVars, nil
}
