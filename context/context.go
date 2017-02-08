package context

import (
	"encoding/json"
	"github.com/polydawn/meep"
	"io/ioutil"
	"path"
)

type ResourceSet struct {
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
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

	err = json.Unmarshal(file, &c)
	if err != nil {
		return nil, meep.New(
			&ContextLoadingError{Filename: filename},
			meep.Cause(err),
		)
	}

	c.BaseDir = path.Dir(filename)

	return &c, nil
}
