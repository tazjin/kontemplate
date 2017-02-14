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

	c.ResourceSets = *flattenResourceSetCollections(&c.ResourceSets)
	c.BaseDir = path.Dir(filename)

	return &c, nil
}

// Flattens resource set collections, i.e. resource sets that themselves have an additional 'include' field set.
// Those will be regarded as a short-hand for including multiple resource sets from a subfolder.
// See https://github.com/tazjin/kontemplate/issues/9 for more information.
func flattenResourceSetCollections(rs *[]ResourceSet) *[]ResourceSet {
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

	return &flattened
}
