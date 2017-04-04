package context

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/polydawn/meep"
	"github.com/tazjin/kontemplate/util"
)

type ResourceSet struct {
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`

	// Fields for resource set collections
	Include []ResourceSet `json:"include"`
	Parent  string
}

type Context struct {
	Name         string                 `json:"context"`
	Global       map[string]interface{} `json:"global"`
	ResourceSets []ResourceSet          `json:"include"`
	BaseDir      string
}

type ContextLoadingError struct {
	meep.AllTraits
	Filename string
}

// Attempt to load and deserialise a Context from the specified file.
func LoadContextFromFile(filename string) (*Context, error) {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, meep.New(
			&ContextLoadingError{Filename: filename},
			meep.Cause(err),
		)
	}

	var c Context

	if strings.HasSuffix(filename, "json") {
		err = json.Unmarshal(file, &c)
	} else if strings.HasSuffix(filename, "yaml") || strings.HasSuffix(filename, "yml") {
		err = yaml.Unmarshal(file, &c)
	} else {
		return nil, meep.New(
			&ContextLoadingError{Filename: filename},
			meep.Cause(fmt.Errorf("File format not supported. Must be JSON or YAML.")),
		)
	}

	if err != nil {
		return nil, meep.New(
			&ContextLoadingError{Filename: filename},
			meep.Cause(err),
		)
	}

	c.ResourceSets = flattenResourceSetCollections(&c.ResourceSets)
	c.BaseDir = path.Dir(filename)
	c.ResourceSets = loadAllDefaultValues(&c)

	return &c, nil
}

// Flattens resource set collections, i.e. resource sets that themselves have an additional 'include' field set.
// Those will be regarded as a short-hand for including multiple resource sets from a subfolder.
// See https://github.com/tazjin/kontemplate/issues/9 for more information.
func flattenResourceSetCollections(rs *[]ResourceSet) []ResourceSet {
	flattened := make([]ResourceSet, 0)

	for _, r := range *rs {
		if len(r.Include) == 0 {
			flattened = append(flattened, r)
		} else {
			for _, subResourceSet := range r.Include {
				subResourceSet.Parent = r.Name
				subResourceSet.Name = path.Join(r.Name, subResourceSet.Name)
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

	// Attempt to load YAML values
	y, err := ioutil.ReadFile(path.Join(c.BaseDir, rs.Name, "default.yaml"))
	if err == nil {
		yaml.Unmarshal(y, &defaultVars)
		return util.Merge(&defaultVars, &rs.Values)
	}

	// Attempt to load JSON values
	j, err := ioutil.ReadFile(path.Join(c.BaseDir, rs.Name, "default.json"))
	if err == nil {
		json.Unmarshal(j, &defaultVars)
		return util.Merge(&defaultVars, &rs.Values)
	}

	// The actual error is not inspected here. The reasoning for this is that in case of serious problems (e.g.
	// permission issues with the folder / folder not existing) failure will occur a bit later anyways.
	// Otherwise we'd have to differentiate between file-not-found-errors (no default values specified) and other
	// errors here.
	return &rs.Values
}
